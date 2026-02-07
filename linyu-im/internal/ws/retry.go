package ws

import "time"

var RetryIntervals = []time.Duration{
	5 * time.Second,
	10 * time.Second,
	30 * time.Second,
}

type RetryMessage struct {
	Key        string
	SeqId      string
	Device     string
	UserId     string
	Data       any
	RetryCount int
}
