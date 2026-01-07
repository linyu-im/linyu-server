package ws

import "sync"

var Manager = NewClientManager()

type ClientManager struct {
	Users map[string]map[string]*Client
	lock  sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		Users: make(map[string]map[string]*Client),
	}
}

func (m *ClientManager) Join(userId string, client *Client) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.Users[userId] == nil {
		m.Users[userId] = make(map[string]*Client)
	}
	if old, ok := m.Users[userId][client.DeviceId]; ok {
		_ = old.Conn.Close()
	}
	m.Users[userId][client.DeviceId] = client
}

func (m *ClientManager) Leave(userId string, deviceId string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	devices, ok := m.Users[userId]
	if !ok {
		return
	}
	if client, ok := devices[deviceId]; ok {
		_ = client.Conn.Close()
		delete(devices, deviceId)
	}

	if len(devices) == 0 {
		delete(m.Users, userId)
	}
}

func (m *ClientManager) SendToUser(userId string, msg []byte) {
	if devices, ok := m.Users[userId]; ok {
		for _, client := range devices {
			select {
			case client.Send <- msg:
			default:
			}
		}
	}
}
