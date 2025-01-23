package channels

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

// UnsubscribeFromChannelResponse is the response of unsubscribing from a channel.
type UnsubscribeFromChannelResponse struct {
	zulip.APIResponseBase
	unsubscribeFromChannelResponseData
}

// unsubscribeFromChannelResponseData represents the JSON data returned by the API
// in the response of unsubscribing from a channel.
type unsubscribeFromChannelResponseData struct {
	NotRemoved []string `json:"not_removed"`
	Removed    []string `json:"removed"`
}

func (s *UnsubscribeFromChannelResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &s.unsubscribeFromChannelResponseData); err != nil {
		return err
	}

	return nil
}

type unsubscribeFromChannelOptions struct {
	principals struct {
		fieldName string
		value     any
	}
}

// UnsubscribeFromChannelOption is the type of the options for unsubscribing from a channel.
type UnsubscribeFromChannelOption func(*unsubscribeFromChannelOptions)

// Principals is the list of users to unsubscribe from the channel.
// A list of user IDs (preferred) or Zulip API email addresses of the users to be subscribed to
// or unsubscribed from the channels specified in the subscriptions parameter. If not provided,
// then the requesting user/bot is unsubscribed.
func Principals[T []int | []string](users T) UnsubscribeFromChannelOption {
	return func(args *unsubscribeFromChannelOptions) {
		args.principals.fieldName = "principals"
		args.principals.value = users
	}
}

// UnsubscribeFromChannel unsubscribes from a channel or a list of channels.
func (svc *Service) UnsubscribeFromChannel(ctx context.Context, subscriptions []string, options ...UnsubscribeFromChannelOption) (*UnsubscribeFromChannelResponse, error) {
	const (
		method = http.MethodDelete
		path   = "/api/v1/users/me/subscriptions"
	)

	channelJSONList, err := json.Marshal(subscriptions)
	if err != nil {
		return nil, err
	}

	msg := map[string]any{
		"subscriptions": string(channelJSONList),
	}

	opts := unsubscribeFromChannelOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.principals.value != nil {
		jsonVal, _ := json.Marshal(opts.principals.value)
		msg[opts.principals.fieldName] = string(jsonVal)
	}

	resp := UnsubscribeFromChannelResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
