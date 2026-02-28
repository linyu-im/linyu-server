package ws

import (
	"sync"
	"time"

	"github.com/RussellLuo/timingwheel"
)

var Manager = NewClientManager()

type Connection interface {
	Close() error
}

type ClientManager struct {
	Users         map[string]map[string]*Client
	lock          sync.RWMutex
	RetryMessages sync.Map
	Tw            *timingwheel.TimingWheel
}

func NewClientManager() *ClientManager {
	tw := timingwheel.NewTimingWheel(time.Second, 60)
	tw.Start()
	m := &ClientManager{
		Tw:    tw,
		Users: make(map[string]map[string]*Client),
	}
	InitTask(m)
	return m
}

func (m *ClientManager) Join(userId string, client *Client) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.Users[userId] == nil {
		m.Users[userId] = make(map[string]*Client)
	}
	if old, ok := m.Users[userId][client.Device]; ok {
		_ = old.Conn.Close()
	}
	m.Users[userId][client.Device] = client
}

func (m *ClientManager) Leave(userId string, device string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.LeaveUnlock(userId, device)
}

func (m *ClientManager) LeaveUnlock(userId string, device string) {
	devices, ok := m.Users[userId]
	if !ok {
		return
	}
	if client, ok := devices[device]; ok {
		_ = client.Conn.Close()
		delete(devices, device)
	}

	if len(devices) == 0 {
		delete(m.Users, userId)
	}
}

func (m *ClientManager) SendToUsers(seqId string, userIds []string, data any) {
	for _, userId := range userIds {
		m.SendToUser(seqId, userId, data)
	}
}

func (m *ClientManager) SendToUser(seqId string, userId string, data any) {
	m.lock.RLock()
	devices, ok := m.Users[userId]
	m.lock.RUnlock()
	if !ok {
		return
	}
	for device, client := range devices {
		m.Send(client, seqId, device, data)
		m.AddRetryMessage(seqId, userId, device, data)
	}
}

func (m *ClientManager) SendToUserOnDevice(seqId string, userId, device string, data any) {
	m.lock.RLock()
	devices, ok := m.Users[userId]
	m.lock.RUnlock()
	if !ok {
		return
	}
	if client, ok := devices[device]; ok {
		m.Send(client, seqId, device, data)
		m.AddRetryMessage(seqId, userId, device, data)
	}
}

func (m *ClientManager) Send(client *Client, seqId, device string, data any) {
	client.SendMsg(SucceedResponse(seqId, device, "server", data))
}

func (m *ClientManager) GetRetryMessageKey(device string, seqId string) string {
	return device + ":" + seqId
}

func (m *ClientManager) AddRetryMessage(seqId, userId, device string, data any) {
	key := m.GetRetryMessageKey(device, seqId)
	if _, loaded := m.RetryMessages.Load(key); loaded {
		return
	}
	retryMsg := &RetryMessage{
		Key:        key,
		SeqId:      seqId,
		UserId:     userId,
		Device:     device,
		Data:       data,
		RetryCount: 0,
	}
	_, loaded := m.RetryMessages.LoadOrStore(key, retryMsg)
	if loaded {
		return
	}
	m.ScheduleRetry(retryMsg)
}

func (m *ClientManager) DeleteRetryMessage(device string, seqId string) {
	m.RetryMessages.Delete(m.GetRetryMessageKey(device, seqId))
}

func (m *ClientManager) ScheduleRetry(retryMsg *RetryMessage) {
	if retryMsg.RetryCount >= len(RetryIntervals) {
		m.RetryMessages.Delete(retryMsg.Key)
		return
	}
	delay := RetryIntervals[retryMsg.RetryCount]
	m.Tw.AfterFunc(delay, func() {
		value, ok := m.RetryMessages.Load(retryMsg.Key)
		if !ok {
			return
		}
		p := value.(*RetryMessage)
		m.SendToUserOnDevice(p.SeqId, p.UserId, p.Device, p.Data)
		p.RetryCount++
		m.ScheduleRetry(p)
	})
}

// CleanExpiredClients 心跳清理
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
