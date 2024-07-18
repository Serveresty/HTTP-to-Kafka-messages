package kafka

import (
	"sync"

	"github.com/IBM/sarama"
)

type Responce struct {
	ResponseChannels map[string]chan *sarama.ConsumerMessage
	mtx              sync.Mutex
}

func NewResponces() *Responce {
	return &Responce{
		ResponseChannels: make(map[string]chan *sarama.ConsumerMessage),
	}
}

func (r *Responce) SetChan(id string, ch chan *sarama.ConsumerMessage) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
}

func (r *Responce) GetChan(id string) chan *sarama.ConsumerMessage {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return nil
}

func (r *Responce) RemoveChan(id string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
}
