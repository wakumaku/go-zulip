package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetUserPresenceAllResponse struct {
	zulip.APIResponseBase
	getUserPresenceAllResponseData
}

type getUserPresenceAllResponseData struct {
	Presences map[string]struct {
		Aggregated struct {
			Status    string `json:"status"` // active, idle
			Timestamp int    `json:"timestamp"`
		} `json:"aggregated"`
		Website struct {
			Status    string `json:"status"`
			Timestamp int    `json:"timestamp"`
		} `json:"website"`
	} `json:"presences"`
	ServerTimestamp float64 `json:"server_timestamp"`
}

func (g *GetUserPresenceAllResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getUserPresenceAllResponseData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) GetUserPresenceAll(ctx context.Context) (*GetUserPresenceAllResponse, error) {
	const (
		path   = "/api/v1/realm/presence"
		method = http.MethodGet
	)

	resp := GetUserPresenceAllResponse{}
	if err := svc.client.DoRequest(ctx, method, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
