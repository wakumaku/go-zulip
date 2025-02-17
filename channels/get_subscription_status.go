package channels

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetSubscriptionStatusResponse struct {
	zulip.APIResponseBase
	getSubscriptionStatusResponseData
}

type getSubscriptionStatusResponseData struct {
	IsSubscribed bool `json:"is_subscribed"`
}

func (g *GetSubscriptionStatusResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getSubscriptionStatusResponseData); err != nil {
		return err
	}

	return nil
}

// GetSubscriptionStatus Check whether a user is subscribed to a channel.
func (svc *Service) GetSubscriptionStatus(ctx context.Context, userID, streamID int) (*GetSubscriptionStatusResponse, error) {
	const (
		method = http.MethodGet
		path   = "/api/v1/users/%d/subscriptions/%d"
	)

	patchPath := fmt.Sprintf(path, userID, streamID)

	resp := GetSubscriptionStatusResponse{}
	if err := svc.client.DoRequest(ctx, method, patchPath, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
