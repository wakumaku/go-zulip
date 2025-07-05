// Package org provides functionality for managing Zulip server and organization settings.
//
// Implemented features:
//   - Upload custom emoji (from file path, bytes, or reader)
//
// See https://zulip.com/api/ for the complete API documentation.
package org

import "github.com/wakumaku/go-zulip"

type Service struct {
	client zulip.RESTClient
}

func NewService(c zulip.RESTClient) *Service {
	return &Service{client: c}
}
