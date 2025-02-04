package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type UpdateUserPresenceResponse struct {
	zulip.APIResponseBase
	updateUserPresenceResponseData
}

type updateUserPresenceResponseData struct {
	PresenceLastUpdateID int `json:"presence_last_update_id"`
	Presences            map[string]struct {
		ActiveTimestamp int `json:"active_timestamp"`
		IdleTimestamp   int `json:"idle_timestamp"`
	} `json:"presences"`
	ServerTimestamp float64 `json:"server_timestamp"`
}

func (g *UpdateUserPresenceResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.updateUserPresenceResponseData); err != nil {
		return err
	}

	return nil
}

type UserPresence string

const (
	UserPresenceActive UserPresence = "active"
	UserPresenceIdle   UserPresence = "idle"
)

func (svc *Service) UpdateUserPresence(ctx context.Context, status UserPresence) (*UpdateUserPresenceResponse, error) {
	const (
		path   = "/api/v1/users/me/presence"
		method = http.MethodPost
	)

	msg := map[string]any{
		"status": status,
	}

	resp := UpdateUserPresenceResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
