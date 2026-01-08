package eventbus

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"sync"
)

type LocalEventBus struct {
	Handlers map[string][]event.Handler
	Lock     sync.RWMutex
}

func NewLocalEventBus() *LocalEventBus {
	return &LocalEventBus{
		Handlers: make(map[string][]event.Handler),
	}
}

func (b *LocalEventBus) Subscribe(e event.Event, handler event.Handler) {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	b.Handlers[e.EventName()] = append(b.Handlers[e.EventName()], handler)
}

func (b *LocalEventBus) Publish(e event.Event) error {
	b.Lock.RLock()
	handlers := b.Handlers[e.EventName()]
	b.Lock.RUnlock()
	for _, h := range handlers {
		go func(handler event.Handler) {
			_ = handler(e)
		}(h)
	}
	return nil
}
