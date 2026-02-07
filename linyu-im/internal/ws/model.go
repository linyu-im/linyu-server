package ws

import "encoding/json"

type Request struct {
	Device string      `json:"device"`
	SeqId  string      `json:"seqId"`          // 消息序列
	Route  string      `json:"route"`          // ws请求路由
	Data   interface{} `json:"data,omitempty"` // 消息内容json
}

type Response struct {
	Device string      `json:"device"`
	SeqId  string      `json:"seqId"`
	Route  string      `json:"route"`
	Code   int         `json:"code"`
	Msg    string      `json:"msg,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func ErrorResponse(seqId, device, route, msg string) []byte {
	response, _ := json.Marshal(&Response{SeqId: seqId, Device: device, Route: route, Code: 500, Msg: msg})
	return response
}

func SucceedResponse(seqId, device, route string, data any) []byte {
	response, _ := json.Marshal(&Response{SeqId: seqId, Device: device, Route: route, Code: 200, Data: data})
	return response
}
