package consumer

import (
	"encoding/json"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"github.com/linyu-im/linyu-server/linyu-im/internal/ws"
)

func WsDataConsumerHandler(e event.Event) error {
	me, ok := e.(event.WsDataEvent)
	if !ok {
		return nil
	}
	msgBytes, _ := json.Marshal(e)
	ws.Manager.SendToUsers(me.ToUserIds, msgBytes)
	return nil
}
