package constant

type wsDataType struct {
	Message string // 消息
	Apply   string // 申请
	Notify  string //通知
}

// WsDataType websocket数据类型
var WsDataType = wsDataType{
	Message: "message",
	Apply:   "apply",
	Notify:  "notify",
}
