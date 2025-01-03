package specialty

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type FetchAPIKeyResponse struct {
	zulip.APIResponseBase
	fetchAPIKeyData
}

type fetchAPIKeyData struct {
	UserID int    `json:"user_id"`
	APIKey string `json:"api_key"`
	Email  string `json:"email"`
}

func (f *FetchAPIKeyResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &f.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &f.fetchAPIKeyData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) FetchAPIKeyProduction(ctx context.Context, username, password string) (*FetchAPIKeyResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/fetch_api_key"
	)

	msg := map[string]any{
		"username": username,
		"password": password,
	}

	resp := FetchAPIKeyResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (svc *Service) FetchAPIKeyDevelopment(ctx context.Context, username string) (*FetchAPIKeyResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/dev_fetch_api_key"
	)

	msg := map[string]any{
		"username": username,
	}

	resp := FetchAPIKeyResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
