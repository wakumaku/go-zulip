package channels

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetChannelIDResponse struct {
	zulip.APIResponseBase
	getChannelIDResponseData
}

type getChannelIDResponseData struct {
	StreamID int `json:"stream_id"`
}

func (g *GetChannelIDResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getChannelIDResponseData); err != nil {
		return err
	}

	return nil
}

// GetChannelID Get the unique ID of a given channel.
func (svc *Service) GetChannelID(ctx context.Context, stream string) (*GetChannelIDResponse, error) {
	const (
		method = http.MethodGet
		path   = "/api/v1/get_stream_id"
	)

	msg := map[string]any{
		"stream": stream,
	}

	resp := GetChannelIDResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
