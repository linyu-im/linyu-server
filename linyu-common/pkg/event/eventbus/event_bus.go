package eventbus

import "github.com/linyu-im/linyu-server/linyu-common/pkg/event"

type EventBus interface {
	Publish(event event.Event) error
	Subscribe(event event.Event, handler event.Handler)
}

var DefaultEventBus EventBus = NewLocalEventBus()
