package kafka

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

func InitProducer() sarama.SyncProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKER")}, cfg)
	if err != nil {
		log.Fatal("failed to start kafka producer: ", err)
	}

	return producer
}

func InitConsumer() sarama.Consumer {
	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumer, err := sarama.NewConsumer([]string{os.Getenv("KAFKA_BROKER")}, cfg)
	if err != nil {
		log.Fatal("failed to start consumer group: ", err)
	}

	return consumer
}
