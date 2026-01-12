package ws

import (
	"encoding/json"
	"sync"
	"time"
)

func InitRoute() {
	Register("heartbeat", HeartbeatHandler)
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
			response, _ := json.Marshal(ErrorResponse(req.Seq, req.Route, err.Error()))
			client.SendMsg(response)
		} else {
			response, _ := json.Marshal(SucceedResponse(req.Seq, req.Route, data))
			client.SendMsg(response)
		}
	} else {
		response, _ := json.Marshal(ErrorResponse(req.Seq, req.Route, "The route does not exist."))
		client.SendMsg(response)
	}
}
