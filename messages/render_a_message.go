package messages

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type RenderAMessageResponse struct {
	zulip.APIResponseBase
	renderAMessageResponseData
}

type renderAMessageResponseData struct {
	Rendered string `json:"rendered"`
}

func (aer *RenderAMessageResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &aer.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &aer.renderAMessageResponseData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) RenderAMessage(ctx context.Context, content string) (*RenderAMessageResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/messages/render"
	)

	msg := map[string]any{
		"content": content,
	}

	resp := RenderAMessageResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
