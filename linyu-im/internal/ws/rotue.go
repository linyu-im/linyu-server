package ws

import (
	"sync"
	"time"
)

func init() {
	Register("heartbeat", HeartbeatHandler)
	Register("ack", AckHandler)
}

type RouteHandler func(client *Client, req *Request) (interface{}, error)

var routes = make(map[string]RouteHandler)
var routeMutex sync.RWMutex

func Register(key string, value RouteHandler) {
	routeMutex.Lock()
	defer routeMutex.Unlock()
	routes[key] = value
}

func GetHandlers(key string) (value RouteHandler, ok bool) {
	routeMutex.RLock()
	defer routeMutex.RUnlock()
	value, ok = routes[key]
	return
}

func ProcessData(client *Client, req *Request) {
	//更新心跳
	currentTime := uint64(time.Now().Unix())
	client.HeartbeatTime = currentTime
	if routeHandler, ok := GetHandlers(req.Route); ok {
		data, err := routeHandler(client, req)
		if err != nil {
			client.SendMsg(ErrorResponse(req.SeqId, req.Device, req.Route, err.Error()))
		} else {
			client.SendMsg(SucceedResponse(req.SeqId, req.Device, req.Route, data))
		}
	} else {
		client.SendMsg(ErrorResponse(req.SeqId, req.Device, req.Route, "The route does not exist."))
	}
}
