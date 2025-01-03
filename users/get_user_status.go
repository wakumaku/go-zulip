package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetUserStatusResponse struct {
	zulip.APIResponseBase
	getUserStatusResponseData
}

type getUserStatusResponseData struct {
	Status struct {
		EmojiCode    string             `json:"emoji_code"`
		EmojiName    string             `json:"emoji_name"`
		ReactionType zulip.ReactionType `json:"reaction_type"`
		StatusText   string             `json:"status_text"`
	} `json:"status"`
}

func (g *GetUserStatusResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getUserStatusResponseData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) GetUserStatus(ctx context.Context, id int) (*GetUserStatusResponse, error) {
	const (
		path   = "/api/v1/users/%d/status"
		method = http.MethodGet
	)
	pathPatch := fmt.Sprintf(path, id)

	resp := GetUserStatusResponse{}
	if err := svc.client.DoRequest(ctx, method, pathPatch, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
