// Package realtime provides real-time event handling for Zulip.
//
// Implemented features:
//   - Register event queue (with various event types and filters)
//   - Get events from event queue (long polling)
//   - Delete event queue
//   - Support for various event types (messages, presence, typing, etc.)
//
// See https://zulip.com/api/ for the complete API documentation.
package realtime

import "github.com/wakumaku/go-zulip"

type Service struct {
	client zulip.RESTClient
}

func NewService(c zulip.RESTClient) *Service {
	return &Service{client: c}
}
