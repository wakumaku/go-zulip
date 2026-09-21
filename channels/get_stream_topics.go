package channels

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type Topic struct {
	MaxID int    `json:"max_id"`
	Name  string `json:"name"`
}

type GetStreamTopicsResponse struct {
	zulip.APIResponseBase
	getStreamTopicsResponseData
}

type getStreamTopicsResponseData struct {
	Topics []Topic `json:"topics"`
}

func (g *GetStreamTopicsResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getStreamTopicsResponseData); err != nil {
		return err
	}

	return nil
}

type getStreamTopicsOptions struct {
	allowEmptyTopicName *bool
}

type GetStreamTopicsOption func(*getStreamTopicsOptions)

// AllowEmptyTopicName tells the server the client can process empty strings
// as topic names. Defaults to false when not set.
func AllowEmptyTopicName(allow bool) GetStreamTopicsOption {
	return func(o *getStreamTopicsOptions) {
		o.allowEmptyTopicName = &allow
	}
}

// GetStreamTopics fetches the topics in a channel, most recent first. See:
// https://zulip.com/api/get-stream-topics
func (svc *Service) GetStreamTopics(ctx context.Context, streamID int, opts ...GetStreamTopicsOption) (*GetStreamTopicsResponse, error) {
	const (
		method = http.MethodGet
		path   = "/api/v1/users/me/%d/topics"
	)

	options := getStreamTopicsOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	msg := map[string]any{}
	if options.allowEmptyTopicName != nil {
		msg["allow_empty_topic_name"] = *options.allowEmptyTopicName
	}

	resp := GetStreamTopicsResponse{}
	if err := svc.client.DoRequest(ctx, method, fmt.Sprintf(path, streamID), msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
