package channels

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetChannelSubscribersResponse struct {
	zulip.APIResponseBase
	getChannelSubscribersResponseData
}

type getChannelSubscribersResponseData struct {
	Subscribers []int `json:"subscribers"`
}

func (g *GetChannelSubscribersResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getChannelSubscribersResponseData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) GetChannelSubscribers(ctx context.Context, streamID int) (*GetChannelSubscribersResponse, error) {
	const (
		method = http.MethodGet
		path   = "/api/v1/streams/%d/members"
	)

	patchPath := fmt.Sprintf(path, streamID)

	resp := GetChannelSubscribersResponse{}
	if err := svc.client.DoRequest(ctx, method, patchPath, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
