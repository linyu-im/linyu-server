package api

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event/eventbus"
	"github.com/linyu-im/linyu-server/linyu-im/internal/consumer"
)

func init() {
	eventbus.DefaultEventBus.Subscribe(event.MessageEvent{}, consumer.MessageConsumerHandler)
}
