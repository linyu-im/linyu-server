package subscriber

import (
	"encoding/json"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event/eventbus"
	"github.com/linyu-im/linyu-server/linyu-im/internal/ws"
)

func init() {
	eventbus.Register[event.WsDataEvent](WsDataEventHandler)
}

func WsDataEventHandler(e *event.WsDataEvent) error {
	msgBytes, _ := json.Marshal(e.Data)
	ws.Manager.SendToUsers(e.ToUserIds, msgBytes)
	return nil
}
