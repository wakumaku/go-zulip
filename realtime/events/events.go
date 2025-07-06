// Package events provides interfaces and types for handling Zulip events.
package events

type EventType string

type Event interface {
	EventID() int
	EventType() EventType
	EventOp() string
}
