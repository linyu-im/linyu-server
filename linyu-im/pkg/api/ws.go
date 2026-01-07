package api

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-im/internal/ws"
)

func init() {
	route.Register("GET", "/ws", ws.WsHandler, false)
}
