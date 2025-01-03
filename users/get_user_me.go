package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetUserMeResponse struct {
	zulip.APIResponseBase
	getUserMeResponseData
}

type getUserMeResponseData struct {
	UserID         int                         `json:"user_id"`
	Role           zulip.OrganizationRoleLevel `json:"role"`
	AvatarUrl      string                      `json:"avatar_url"`
	AvatarVersion  int                         `json:"avatar_version"`
	DateJoined     string                      `json:"date_joined"`
	DeliveryEmail  string                      `json:"delivery_email"`
	Email          string                      `json:"email"`
	FullName       string                      `json:"full_name"`
	IsActive       bool                        `json:"is_active"`
	IsAdmin        bool                        `json:"is_admin"`
	IsBillingAdmin bool                        `json:"is_billing_admin"`
	IsBot          bool                        `json:"is_bot"`
	IsGuest        bool                        `json:"is_guest"`
	IsOwner        bool                        `json:"is_owner"`
	MaxMessageId   int                         `json:"max_message_id"`
	ProfileData    map[string]struct {
		Value         string `json:"value"`
		RenderedValue string `json:"rendered_value,omitempty"`
	} `json:"profile_data"`
}

func (aer *GetUserMeResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &aer.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &aer.getUserMeResponseData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) GetUserMe(ctx context.Context) (*GetUserMeResponse, error) {
	const (
		path   = "/api/v1/users/me"
		method = http.MethodGet
	)

	msg := map[string]any{}

	resp := GetUserMeResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
