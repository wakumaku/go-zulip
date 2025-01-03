package messages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/wakumaku/go-zulip"
)

type SendMessageType string

const (
	ToChannel SendMessageType = "channel"
	ToDirect  SendMessageType = "direct"
	ToStream  SendMessageType = "stream"
	// ToPrivate
	// Deprecated: In Zulip 7.0 (feature level 174), "direct" was added as the
	// preferred way to request a direct message, deprecating the original
	// "private". While "private" is still supported for requesting direct
	// messages, clients are encouraged to use to the modern convention with
	// servers that support it, because support for "private" will eventually
	// be removed.
	ToPrivate SendMessageType = "private"
)

// SendMessageRecipient For channel messages, either the name or integer ID of the
// channel. For direct messages, either a list containing integer user IDs or
// a list containing string Zulip API email addresses.
type SendMessageRecipient interface {
	Recipient() any
	SendMessageType() SendMessageType
}

type ToChannelID int

func (t ToChannelID) Recipient() any {
	return t
}

func (t ToChannelID) SendMessageType() SendMessageType {
	return ToChannel
}

type ToChannelName string

func (t ToChannelName) Recipient() any {
	return t
}

func (t ToChannelName) SendMessageType() SendMessageType {
	return ToChannel
}

type toChannelTopic struct {
	channelNameID SendMessageRecipient
	toTopic       string
}

func (t toChannelTopic) Recipient() any {
	return t.channelNameID.Recipient()
}

func (t toChannelTopic) SendMessageType() SendMessageType {
	return ToChannel
}

func (t toChannelTopic) Topic() string {
	return t.toTopic
}

func ToChannelTopic(channelNameID SendMessageRecipient, topic string) SendMessageRecipient {
	return toChannelTopic{
		channelNameID: channelNameID,
		toTopic:       topic,
	}
}

type ToUserID int

func (t ToUserID) Recipient() any {
	return t
}

func (t ToUserID) SendMessageType() SendMessageType {
	return ToDirect
}

type ToUserIDs []int

func (t ToUserIDs) Recipient() any {
	return t
}

func (t ToUserIDs) SendMessageType() SendMessageType {
	return ToDirect
}

type ToUserName string

func (t ToUserName) Recipient() any {
	return t
}

func (t ToUserName) SendMessageType() SendMessageType {
	return ToDirect
}

type ToUserNames []string

func (t ToUserNames) Recipient() any {
	return t
}

func (t ToUserNames) SendMessageType() SendMessageType {
	return ToDirect
}

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

func (svc *Service) SendMessage(ctx context.Context, to SendMessageRecipient, content string, options ...SendMessageOption) (*SendMessageResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/messages"
	)

	var toRecipient any
	switch t := to.(type) {
	case ToUserIDs, ToUserNames:
		v, err := json.Marshal(t.Recipient())
		if err != nil {
			return nil, err
		}
		toRecipient = string(v)
	case toChannelTopic:
		toRecipient = t.Recipient()
		options = append(options, ToTopic(t.Topic()))
	default:
		toRecipient = to.Recipient()
	}

	msg := map[string]any{
		"to":      toRecipient,
		"type":    to.SendMessageType(),
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
