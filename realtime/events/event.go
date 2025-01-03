package events

type EventType string

type Event interface {
	EventID() int
	EventType() EventType
	EventOp() string
}
