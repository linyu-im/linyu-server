package eventbus

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
)

type Factory func() event.Event

type EventBus interface {
	Publish(event event.Event) error
	Subscribe(eventName string, factory Factory, handler event.Handler)
}

type subscription struct {
	eventName string
	factory   Factory
	handler   event.Handler
}

var handles []subscription

func Register[T any, PT interface {
	*T
	event.Event
}](handler func(PT) error) {
	temp := PT(new(T))
	eventName := temp.EventName()
	factory := func() event.Event { return PT(new(T)) }

	wrappedHandler := func(e event.Event) error {
		if typedEvent, ok := e.(PT); ok {
			return handler(typedEvent)
		}
		if typedEvent, ok := e.(T); ok {
			return handler(&typedEvent)
		}
		return nil
	}

	handles = append(handles, subscription{
		eventName: eventName,
		factory:   factory,
		handler:   wrappedHandler,
	})
}

var GlobalBus EventBus

func InitEventBus() EventBus {
	if config.C.EventBus.Type == "" {
		panic("event bus type not set")
	}
	var bus EventBus
	var err error

	switch config.C.EventBus.Type {
	case config.RocketEventBusType:
		bus, err = NewRocketMQEventBus()
		if err != nil {
			panic("rocketmq init error: " + err.Error())
		}
	case config.LocalEventBusType:
		bus = NewLocalEventBus()
	default:
		panic("event bus type not supported: " + config.C.EventBus.Type)
	}

	for _, s := range handles {
		bus.Subscribe(s.eventName, s.factory, s.handler)
	}
	handles = nil

	GlobalBus = bus
	return bus
}
