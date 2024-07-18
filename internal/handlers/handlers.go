package handlers

import (
	"Messagio/internal/kafka"
	"Messagio/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	DB               *sql.DB
	Producer         sarama.SyncProducer
	ResponceChannels *kafka.Responce
}

func SetupRoutes(r *mux.Router, db *sql.DB, prod sarama.SyncProducer, respCh *kafka.Responce) {
	h := &Handler{DB: db, Producer: prod, ResponceChannels: respCh}
	r.HandleFunc("/messages", h.CreateMessage).Methods("POST")
	r.HandleFunc("/messages/stats", h.GetStats).Methods("GET")
}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()

	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg.CreatedAt = time.Now()

	_, err := h.DB.Exec("INSERT INTO messages (data, created_at) VALUES ($1, $2)", msg.Data, msg.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _, err = h.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "messages",
		Value: sarama.StringEncoder(msg.Data),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseCh := make(chan *sarama.ConsumerMessage)
	h.ResponceChannels.SetChan(requestID, responseCh)

	select {
	case responseMsg := <-responseCh:
		_, err := h.DB.Exec("UPDATE messages SET processed = true WHERE data = $1 and created_at = $2", msg.Data, msg.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(responseMsg.Value)
		w.WriteHeader(http.StatusOK)
	case <-time.After(10 * time.Second):
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT data, created_at FROM messages WHERE processed = true")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.Data, &msg.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		msgs = append(msgs, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}
