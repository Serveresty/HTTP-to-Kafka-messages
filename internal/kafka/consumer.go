package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

func ReadMessages(partConsumer sarama.PartitionConsumer, respCh *Responce) {
	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				log.Println("Channel closed, exiting goroutine")
				return
			}
			respID := string(msg.Key)
			respCh.mtx.Lock()
			ch, exists := respCh.ResponseChannels[respID]
			if exists {
				ch <- msg
				delete(respCh.ResponseChannels, respID)
			}
			respCh.mtx.Unlock()
		}
	}
}
