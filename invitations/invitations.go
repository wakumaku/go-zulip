// Package invitations provides functionality for managing Zulip invitations.
//
// Implemented features:
//   - Create reusable invitation links
//
// See https://zulip.com/api/ for the complete API documentation.
package invitations

import "github.com/wakumaku/go-zulip"

type Service struct {
	client zulip.RESTClient
}

func NewService(c zulip.RESTClient) *Service {
	return &Service{client: c}
}
