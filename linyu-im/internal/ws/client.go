package ws

import (
	"encoding/json"
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
	send          chan []byte
	Device        string
	HeartbeatTime uint64
	ConnectTime   uint64
}

func NewClient(conn *websocket.Conn, userId string, device string) (client *Client) {
	currentTime := uint64(time.Now().Unix())
	client = &Client{
		UserId:        userId,
		Device:        device,
		IP:            conn.RemoteAddr().String(),
		Conn:          conn,
		send:          make(chan []byte, 100),
		HeartbeatTime: currentTime,
		ConnectTime:   currentTime,
	}
	return
}

func (c *Client) Read() {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("ws read error", zap.String("stack", string(debug.Stack())), zap.Any("r", r))
		}
	}()
	defer func() {
		close(c.send)
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}
		request := &Request{}
		if err := json.Unmarshal(message, request); err != nil {
			c.SendMsg(ErrorResponse("", "", "", "Data formatting error"))
			continue
		}
		ProcessData(c, request)
	}
}

func (c *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("ws write error", zap.String("stack", string(debug.Stack())), zap.Any("r", r))
		}
	}()
	defer func() {
		Manager.Leave(c.UserId, c.Device)
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			_ = c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Info("SendMsg :", zap.String("message", string(debug.Stack())))
		}
	}()
	c.send <- msg
}
