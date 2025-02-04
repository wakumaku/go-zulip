package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetUserPresenceResponse struct {
	zulip.APIResponseBase
	getUserPresenceResponseData
}

type getUserPresenceResponseData struct {
	Presence struct {
		Aggregated struct {
			Status    string `json:"status"` // active, idle
			Timestamp int    `json:"timestamp"`
		} `json:"aggregated"`
		Website struct {
			Status    string `json:"status"`
			Timestamp int    `json:"timestamp"`
		} `json:"website"`
	} `json:"presence"`
}

func (g *GetUserPresenceResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getUserPresenceResponseData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) GetUserPresence(ctx context.Context, userIDorEmail string) (*GetUserPresenceResponse, error) {
	const (
		path   = "/api/v1/users/%s/presence"
		method = http.MethodGet
	)
	pathPatch := fmt.Sprintf(path, userIDorEmail)

	resp := GetUserPresenceResponse{}
	if err := svc.client.DoRequest(ctx, method, pathPatch, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
