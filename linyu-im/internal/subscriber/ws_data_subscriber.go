package subscriber

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event/eventbus"
	"github.com/linyu-im/linyu-server/linyu-im/internal/ws"
)

func init() {
	eventbus.Register[event.WsDataEvent](WsDataEventHandler)
}

func WsDataEventHandler(e *event.WsDataEvent) error {
	ws.Manager.SendToUsers(e.GetSequenceId(), e.ToUserIds, e.Data)
	return nil
}
