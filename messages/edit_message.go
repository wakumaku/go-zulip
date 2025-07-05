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

type PropagateMode string

const (
	// PropagateModeLater The target message and all following messages.
	PropagateModeLater PropagateMode = "change_later"
	// PropagateModeOne Only the target message.
	PropagateModeOne PropagateMode = "change_one"
	// PropagateModeAll All messages in this topic.
	PropagateModeAll PropagateMode = "change_all"
)

type editMessageOptions struct {
	topic struct {
		fieldName string
		value     *string
	}
	propagateMode struct {
		fieldName string
		value     *PropagateMode
	}
	sendNotificationToOldThread struct {
		fieldName string
		value     *bool
	}
	sendNotificationToNewThread struct {
		fieldName string
		value     *bool
	}
	content struct {
		fieldName string
		value     *string
	}
	streamID struct {
		fieldName string
		value     *int
	}
}

type EditMessageOption func(*editMessageOptions) error

// MoveToTopic The topic to move the message(s) to, to request changing the topic.
func MoveToTopic(name string) EditMessageOption {
	return func(o *editMessageOptions) error {
		if strings.TrimSpace(name) == "" {
			return errors.New("topic 'name' is empty")
		}

		o.topic.fieldName = "topic"
		o.topic.value = &name

		return nil
	}
}

// SetPropagateMode Which message(s) should be edited
func SetPropagateMode(change PropagateMode) EditMessageOption {
	return func(o *editMessageOptions) error {
		o.propagateMode.fieldName = "propagate_mode"
		o.propagateMode.value = &change

		return nil
	}
}

// SendNotificationToOldThread Whether to send an automated message to the old topic to notify users where the messages were moved to.
func SendNotificationToOldThread(yes bool) EditMessageOption {
	return func(o *editMessageOptions) error {
		o.sendNotificationToOldThread.fieldName = "send_notification_to_old_thread"
		o.sendNotificationToOldThread.value = &yes

		return nil
	}
}

// SendNotificationToNewThread Whether to send an automated message to the new topic to notify users where the messages came from.
func SendNotificationToNewThread(yes bool) EditMessageOption {
	return func(o *editMessageOptions) error {
		o.sendNotificationToNewThread.fieldName = "send_notification_to_new_thread"
		o.sendNotificationToNewThread.value = &yes

		return nil
	}
}

// NewContent The updated content of the target message.
func NewContent(content string) EditMessageOption {
	return func(o *editMessageOptions) error {
		o.content.fieldName = "content"
		o.content.value = &content

		return nil
	}
}

// SetStreamID The channel ID to move the message(s) to, to request moving messages to another channel.
func SetStreamID(id int) EditMessageOption {
	return func(o *editMessageOptions) error {
		o.streamID.fieldName = "stream_id"
		o.streamID.value = &id

		return nil
	}
}

type EditMessageResponse struct {
	zulip.APIResponseBase
	editMessageResponseData
}

type editMessageResponseData struct {
	DetachedUploads []struct {
		CreateTime int    `json:"create_time"`
		ID         int    `json:"id"`
		Messages   []any  `json:"messages"`
		Name       string `json:"name"`
		PathId     string `json:"path_id"`
		Size       int    `json:"size"`
	} `json:"detached_uploads"`
}

func (e *EditMessageResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &e.editMessageResponseData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) EditMessage(ctx context.Context, id int, options ...EditMessageOption) (*EditMessageResponse, error) {
	const (
		method = http.MethodPatch
		path   = "/api/v1/messages"
	)

	patchPath := fmt.Sprintf("%s/%d", path, id)

	msg := map[string]any{}

	opts := editMessageOptions{}
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}

	if opts.topic.value != nil {
		msg[opts.topic.fieldName] = *opts.topic.value
	}

	if opts.propagateMode.value != nil {
		msg[opts.propagateMode.fieldName] = *opts.propagateMode.value
	}

	if opts.sendNotificationToOldThread.value != nil {
		msg[opts.sendNotificationToOldThread.fieldName] = *opts.sendNotificationToOldThread.value
	}

	if opts.sendNotificationToNewThread.value != nil {
		msg[opts.sendNotificationToNewThread.fieldName] = *opts.sendNotificationToNewThread.value
	}

	if opts.content.value != nil {
		msg[opts.content.fieldName] = *opts.content.value
	}

	if opts.streamID.value != nil {
		msg[opts.streamID.fieldName] = *opts.streamID.value
	}

	resp := EditMessageResponse{}
	if err := svc.client.DoRequest(ctx, method, patchPath, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
