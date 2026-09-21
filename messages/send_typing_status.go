package messages

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

// TypingOp indicates whether typing has started or stopped.
type TypingOp string

const (
	TypingStart TypingOp = "start"
	TypingStop  TypingOp = "stop"
)

// TypingType indicates the kind of message being composed.
type TypingType string

const (
	TypingDirect  TypingType = "direct"
	TypingChannel TypingType = "channel"
)

type SendTypingStatusResponse struct {
	zulip.APIResponseBase
}

func (r *SendTypingStatusResponse) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.APIResponseBase)
}

type sendTypingStatusOptions struct {
	typ       TypingType
	directTo  []int
	channelID *int
	topic     *string
}

type SendTypingStatusOption func(*sendTypingStatusOptions)

// ToDirect marks this as direct-message typing status, addressed to the given
// recipient user IDs. Required when the type is direct.
func ToDirect(userIDs ...int) SendTypingStatusOption {
	return func(o *sendTypingStatusOptions) {
		o.typ = TypingDirect
		o.directTo = userIDs
	}
}

// ToChannelTopic marks this as channel-message typing status, addressed to the
// given channel and topic. Required when the type is channel.
func ToChannelTopic(channelID int, topic string) SendTypingStatusOption {
	return func(o *sendTypingStatusOptions) {
		o.typ = TypingChannel
		o.channelID = &channelID
		o.topic = &topic
	}
}

// SendTypingStatus notifies other clients that the authenticated user has
// started or stopped composing a message. See:
// https://zulip.com/api/set-typing-status
func (svc *Service) SendTypingStatus(ctx context.Context, op TypingOp, opts ...SendTypingStatusOption) (*SendTypingStatusResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/typing"
	)

	options := sendTypingStatusOptions{typ: TypingDirect}
	for _, opt := range opts {
		opt(&options)
	}

	msg := map[string]any{
		"op":   string(op),
		"type": string(options.typ),
	}

	switch options.typ {
	case TypingDirect:
		if len(options.directTo) == 0 {
			return nil, errors.New("go-zulip: SendTypingStatus: ToDirect requires at least one user ID")
		}

		to, err := json.Marshal(options.directTo)
		if err != nil {
			return nil, err
		}

		msg["to"] = string(to)
	case TypingChannel:
		if options.channelID == nil || options.topic == nil {
			return nil, errors.New("go-zulip: SendTypingStatus: ToChannelTopic requires a channel ID and topic")
		}

		msg["stream_id"] = *options.channelID
		msg["topic"] = *options.topic
	}

	resp := SendTypingStatusResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
