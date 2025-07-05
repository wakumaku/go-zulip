// Package messages provides functionality for managing Zulip messages.
//
// Implemented features:
//   - Send messages (to channels/topics or direct messages)
//   - Upload files (from file path, bytes, or reader)
//   - Edit messages
//   - Delete messages
//   - Get messages (with various filters)
//   - Add emoji reactions
//   - Remove emoji reactions
//   - Render messages
//   - Fetch single message
//   - Update personal message flags
//   - Update personal message flags for narrow
//   - Get message read receipts
//
// See https://zulip.com/api/ for the complete API documentation.
package messages

import "github.com/wakumaku/go-zulip"

type Service struct {
	client zulip.RESTClient
}

func NewService(c zulip.RESTClient) *Service {
	return &Service{client: c}
}
