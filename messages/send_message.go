package messages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/messages/recipient"
)

// "direct" for a direct message and "stream" or "channel" for a channel message.
const (
	// In Zulip 9.0 (feature level 248), "channel" was added as an additional value for this parameter to request a channel message.
	toChannel string = "channel"
	// toStream  string = "stream"
	// Direct messages are also known as private messages.
	toDirect string = "direct"
	// ToPrivate
	// Deprecated: In Zulip 7.0 (feature level 174), "direct" was added as the
	// preferred way to request a direct message, deprecating the original
	// "private". While "private" is still supported for requesting direct
	// messages, clients are encouraged to use to the modern convention with
	// servers that support it, because support for "private" will eventually
	// be removed.
	// toPrivate string = "private"
)

type sendMessageOptions struct {
	topic struct {
		fieldName string
		value     *string
	}
	// queueID      string
	// localID      string
	readBySender struct {
		fieldName string
		value     *bool
	}
}

type SendMessageOption func(*sendMessageOptions) error

// ToTopic The topic of the message. Only required for channel
// messages ("type": "stream" or "type": "channel"), ignored otherwise.
func ToTopic(name string) SendMessageOption {
	return func(o *sendMessageOptions) error {
		if strings.TrimSpace(name) == "" {
			return errors.New("topic 'name' is empty")
		}

		o.topic.fieldName = "topic"
		o.topic.value = &name

		return nil
	}
}

// ReadBySender Whether the message should be initially marked read by its sender. If
// unspecified, the server uses a heuristic based on the client name.
func ReadBySender(asRead bool) SendMessageOption {
	return func(o *sendMessageOptions) error {
		o.readBySender.fieldName = "read_by_sender"
		o.readBySender.value = &asRead

		return nil
	}
}

type SendMessageResponse struct {
	zulip.APIResponseBase
	sendMessageResponseData
}

type sendMessageResponseData struct {
	ID                           int `json:"id"`
	AutomaticNewVisibilityPolicy int `json:"automatic_new_visibility_policy"`
}

func (s *SendMessageResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &s.sendMessageResponseData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) SendMessageToChannelTopic(ctx context.Context, channel recipient.Channel, topic string, content string, options ...SendMessageOption) (*SendMessageResponse, error) {
	return svc.SendMessage(ctx, channel, content, append(options, ToTopic(topic))...)
}

func (svc *Service) SendMessageToUsers(ctx context.Context, users recipient.Direct, content string, options ...SendMessageOption) (*SendMessageResponse, error) {
	return svc.SendMessage(ctx, users, content, options...)
}

func (svc *Service) SendMessage(ctx context.Context, to recipient.Recipient, content string, options ...SendMessageOption) (*SendMessageResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/messages"
	)

	var (
		toRecipient   any
		recipientType string
	)

	switch t := to.(type) {
	case recipient.Direct:
		v, err := json.Marshal(t.To())
		if err != nil {
			return nil, err
		}

		toRecipient = string(v)
		recipientType = toDirect
	case recipient.Channel:
		toRecipient = to.To()
		recipientType = toChannel
	default:
		return nil, fmt.Errorf("unsupported recipient type: %T", to)
	}

	msg := map[string]any{
		"to":      toRecipient,
		"type":    recipientType,
		"content": content,
	}

	opts := sendMessageOptions{}
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}

	if opts.topic.value != nil && *opts.topic.value != "" {
		msg[opts.topic.fieldName] = *opts.topic.value
	}

	if opts.readBySender.value != nil {
		msg[opts.readBySender.fieldName] = *opts.readBySender.value
	}

	resp := SendMessageResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
