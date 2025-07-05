package messages

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type DeleteMessageResponse struct {
	zulip.APIResponseBase
}

func (svc *Service) DeleteMessage(ctx context.Context, id int) (*DeleteMessageResponse, error) {
	const (
		method = http.MethodDelete
		path   = "/api/v1/messages"
	)

	patchPath := fmt.Sprintf("%s/%d", path, id)

	resp := DeleteMessageResponse{}
	if err := svc.client.DoRequest(ctx, method, patchPath, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
