// Package specialty provides access to Zulip's specialty endpoints.
//
// Implemented features:
//   - Fetch API key (production and development)
//
// See https://zulip.com/api/ for the complete API documentation.
package specialty

import "github.com/wakumaku/go-zulip"

type Service struct {
	client zulip.RESTClient
}

func NewService(c zulip.RESTClient) *Service {
	return &Service{client: c}
}
