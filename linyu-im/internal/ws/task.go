package ws

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/timer"
	"time"
)

func InitTask(m *ClientManager) {
	timer.Timer(30*time.Second, CleanClientTask, m)
}

func CleanClientTask(param interface{}) bool {
	mgr, ok := param.(*ClientManager)
	if !ok || mgr == nil {
		return true
	}
	return mgr.CleanExpiredClients()
}
