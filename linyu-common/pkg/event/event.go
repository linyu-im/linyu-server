package event

type Event interface {
	EventName() string
	GetSequenceId() string
}

type Handler func(e Event) error
