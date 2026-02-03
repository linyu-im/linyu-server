package constant

type eventType struct {
	WsDataEvent string //ws数据事件
}

// EventType 事件类型
var EventType = eventType{
	WsDataEvent: "WsDataEvent",
}
