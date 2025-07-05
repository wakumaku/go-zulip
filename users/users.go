// Package users provides functionality for managing Zulip users.
//
// Implemented features:
//   - Get a user (by ID or email)
//   - Get own user (me)
//   - Get all users
//   - Create a user
//   - Update a user (by ID or email)
//   - Get user status
//   - Update user status
//   - Get user presence (individual or all users)
//   - Update user presence
//   - Update user settings
//
// See https://zulip.com/api/ for the complete API documentation.
package users

import "github.com/wakumaku/go-zulip"

type Service struct {
	client zulip.RESTClient
}

func NewService(c zulip.RESTClient) *Service {
	return &Service{client: c}
}
