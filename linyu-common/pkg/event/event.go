package event

type Event interface {
	EventName() string
}

type Handler func(e Event) error
