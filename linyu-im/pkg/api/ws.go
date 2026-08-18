package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	_ "github.com/linyu-im/linyu-server/linyu-im/internal/subscriber"
	"github.com/linyu-im/linyu-server/linyu-im/internal/ws"
)

func init() {
	route.Register("GET", "/ws", WsGinApiHandler)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsGinApiHandler(c *gin.Context) {
	userId := c.GetString("userId")
	device := c.GetHeader("device")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := ws.NewClient(conn, userId, device)
	ws.Manager.Join(userId, client)
	go client.Read()
	go client.Write()
}
