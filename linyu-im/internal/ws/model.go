package ws

type Request struct {
	Seq   string      `json:"seq"`            // 消息序列
	Route string      `json:"route"`          // ws请求路由
	Data  interface{} `json:"data,omitempty"` // 消息内容json
}

type Response struct {
	Seq   string      `json:"seq"`
	Route string      `json:"route"`
	Code  int         `json:"code"`
	Msg   string      `json:"msg,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func ErrorResponse(seq, route, msg string) *Response {
	return &Response{Seq: seq, Route: route, Code: 500, Msg: msg}
}

func SucceedResponse(seq, route string, data interface{}) *Response {
	return &Response{Seq: seq, Route: route, Code: 200, Data: data}
}
