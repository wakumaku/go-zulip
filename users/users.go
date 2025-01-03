package users

import "github.com/wakumaku/go-zulip"

type Service struct {
	client zulip.RESTClient
}

func NewService(c zulip.RESTClient) *Service {
	return &Service{client: c}
}
