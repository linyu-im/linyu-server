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
	Send          chan []byte
	DeviceId      string
	HeartbeatTime uint64
	ConnectTime   uint64
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
		close(c.Send)
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			response, _ := json.Marshal(ErrorResponse("", "", err.Error()))
			c.SendMsg(response)
			continue
		}
		request := &Request{}
		if err := json.Unmarshal(message, request); err != nil {
			response, _ := json.Marshal(ErrorResponse("", "", "Data formatting error"))
			c.SendMsg(response)
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

func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Info("SendMsg :", zap.String("message", string(debug.Stack())))
		}
	}()
	c.Send <- msg
}
