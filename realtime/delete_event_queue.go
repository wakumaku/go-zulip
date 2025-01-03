package realtime

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type DeleteEventQueueResponse struct {
	zulip.APIResponseBase
	deleteEventQueueData
}

type deleteEventQueueData struct {
	QueueId string `json:"queue_id"`
}

func (d *DeleteEventQueueResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &d.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &d.deleteEventQueueData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) DeleteEvetQueue(ctx context.Context, queueID string) (*DeleteEventQueueResponse, error) {
	const (
		method = http.MethodDelete
		path   = "/api/v1/events"
	)

	msg := map[string]any{
		"queue_id": queueID,
	}

	resp := DeleteEventQueueResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
