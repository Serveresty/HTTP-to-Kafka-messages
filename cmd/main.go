package main

import (
	"Messagio/internal/database"
	"Messagio/internal/handlers"
	"Messagio/internal/kafka"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {
	db := database.InitDB()
	defer db.Close()

	producer := kafka.InitProducer()
	defer producer.Close()

	consumer := kafka.InitConsumer()
	defer consumer.Close()

	partConsumer, err := consumer.ConsumePartition("messages", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("failed to consume partition: ", err)
	}
	defer partConsumer.Close()

	respCh := kafka.NewResponces()

	go kafka.ReadMessages(partConsumer, respCh)

	r := mux.NewRouter()
	handlers.SetupRoutes(r, db, producer, respCh)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
}
