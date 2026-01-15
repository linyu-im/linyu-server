package ws

import (
	"sync"
	"time"
)

var Manager = NewClientManager()

type ClientManager struct {
	Users map[string]map[string]*Client
	lock  sync.RWMutex
}

func NewClientManager() *ClientManager {
	m := &ClientManager{
		Users: make(map[string]map[string]*Client),
	}
	InitTask(m)
	InitRoute()
	return m
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

	m.LeaveUnlock(userId, deviceId)
}

func (m *ClientManager) LeaveUnlock(userId string, deviceId string) {

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

func (m *ClientManager) SendToUsers(userIds []string, msg []byte) {
	for _, userId := range userIds {
		m.SendToUser(userId, msg)
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

func (m *ClientManager) CleanExpiredClients() bool {
	const timeoutSeconds = 10 * 60

	now := uint64(time.Now().Unix())
	expireThreshold := now - timeoutSeconds

	var expiredClients []struct {
		userId   string
		deviceId string
	}

	m.lock.RLock()
	for userId, devices := range m.Users {
		for deviceId, client := range devices {
			if client.HeartbeatTime < expireThreshold {
				expiredClients = append(expiredClients, struct {
					userId   string
					deviceId string
				}{userId, deviceId})
			}
		}
	}
	m.lock.RUnlock()

	m.lock.Lock()
	defer m.lock.Unlock()

	for _, ec := range expiredClients {
		m.LeaveUnlock(ec.userId, ec.deviceId)
	}
	return true
}
