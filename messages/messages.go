package messages

import "github.com/wakumaku/go-zulip"

type Service struct {
	client zulip.RESTClient
}

func NewService(client zulip.RESTClient) *Service {
	return &Service{client: client}
}
