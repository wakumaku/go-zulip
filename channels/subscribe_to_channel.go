package channels

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

// SubscribeToChannelResponse is the response of subscribing to a channel.
type SubscribeToChannelResponse struct {
	zulip.APIResponseBase
	subscribeToChannelResponseData
}

// subscribeToChannelResponseData represents the JSON data returned by the API
// in the response of subscribing to a channel.
type subscribeToChannelResponseData struct {
	Subscribed        map[string][]string `json:"subscribed"`
	AlreadySubscribed map[string][]string `json:"already_subscribed"`
}

func (aer *SubscribeToChannelResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &aer.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &aer.subscribeToChannelResponseData); err != nil {
		return err
	}

	return nil
}

type subscribeToChannelOptions struct{}

// SubscribeToChannelOption is the type of the options for subscribing to a channel.
type SubscribeToChannelOption func(*subscribeToChannelOptions)

// SubscribeTo is the type of the channel to subscribe to.
type SubscribeTo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// SubscribeToChannel subscribes to (or creates) a channel.
func (svc *Service) SubscribeToChannel(ctx context.Context, list []SubscribeTo, options ...SubscribeToChannelOption) (*SubscribeToChannelResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/users/me/subscriptions"
	)

	channelJSONList, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	msg := map[string]any{
		"subscriptions": string(channelJSONList),
	}

	opts := subscribeToChannelOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	resp := SubscribeToChannelResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
