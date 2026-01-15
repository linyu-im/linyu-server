package event

type WsDataEvent struct {
	FromUserId string      `json:"-"`
	ToUserIds  []string    `json:"-"`
	Type       string      `json:"type"`
	Content    interface{} `json:"content"`
}

func (m WsDataEvent) EventName() string {
	return "WsDataEvent"
}
