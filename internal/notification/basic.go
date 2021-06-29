package notification

import (
	"sync"

	"github.com/Frosin/setserver/internal/api/gen"
)

type BasicNotificator struct {
	consumers []gen.Api_SubscribeServer
	mtx       sync.RWMutex
}

func NewBasicNotificator() Notificator {
	return &BasicNotificator{
		consumers: []gen.Api_SubscribeServer{},
	}
}

func (bn *BasicNotificator) RegisterConsumer(newConsumer gen.Api_SubscribeServer) error {
	bn.mtx.Lock()
	bn.consumers = append(bn.consumers, newConsumer)
	bn.mtx.Unlock()
	return nil
}

func (bn *BasicNotificator) Broadcast(operation, key, value string) error {
	bn.mtx.RLock()
	for _, consumer := range bn.consumers {
		if err := consumer.Send(&gen.Message{
			Name:      key,
			Operation: operation,
			Value:     value,
		}); err != nil {
			bn.mtx.RUnlock()
			return err
		}
	}
	bn.mtx.RUnlock()
	return nil
}
