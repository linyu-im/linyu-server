package ws

import (
	"github.com/gorilla/websocket"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/logger"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

type Client struct {
	UserId        string
	IP            string
	Conn          *websocket.Conn
	Send          chan []byte
	DeviceId      string
	HeartbeatTime uint64
	LoginTime     uint64
}

func NewClient(conn *websocket.Conn, userId string, deviceId string) (client *Client) {
	currentTime := uint64(time.Now().Unix())
	client = &Client{
		UserId:        userId,
		DeviceId:      deviceId,
		IP:            conn.RemoteAddr().String(),
		Conn:          conn,
		Send:          make(chan []byte, 100),
		HeartbeatTime: currentTime,
	}
	return
}

func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("ws read error", zap.String("stack", string(debug.Stack())), zap.Any("r", r))
		}
	}()
	defer func() {
		close(c.Send)
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}
		logger.Log.Info("receive :", zap.String("message", string(message)))
	}
}

func (c *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("ws write error", zap.String("stack", string(debug.Stack())), zap.Any("r", r))
		}
	}()
	defer func() {
		Manager.Leave(c.UserId, c.DeviceId)
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				return
			}
			_ = c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
