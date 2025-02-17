package channels

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetChannelByIDResponse struct {
	zulip.APIResponseBase
	getChannelByIDResponseData
}

type getChannelByIDResponseData struct {
	Stream ChannelInfo `json:"stream"`
}

func (g *GetChannelByIDResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getChannelByIDResponseData); err != nil {
		return err
	}

	return nil
}

// GetChannelByID Fetch details for the channel with the ID.
func (svc *Service) GetChannelByID(ctx context.Context, streamID int) (*GetChannelByIDResponse, error) {
	const (
		method = http.MethodGet
		path   = "/api/v1/streams/%d"
	)

	patchPath := fmt.Sprintf(path, streamID)

	resp := GetChannelByIDResponse{}
	if err := svc.client.DoRequest(ctx, method, patchPath, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
