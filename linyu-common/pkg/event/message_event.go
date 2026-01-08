package event

type MessageEvent struct {
	MessageID  string
	FromUserId string
	ToUserId   string
	Content    string
}

func (m MessageEvent) EventName() string {
	return "MessageEvent"
}
