package consumer

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"github.com/linyu-im/linyu-server/linyu-im/internal/ws"
)

func MessageConsumerHandler(e event.Event) error {
	me, ok := e.(event.MessageEvent)
	if !ok {
		return nil
	}
	ws.Manager.SendToUser(me.ToUserId, []byte(me.Content))
	return nil
}
