package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetUsersResponse struct {
	zulip.APIResponseBase
	getUsersResponseData
}

type getUsersResponseData struct {
	Members []struct {
		UserID         int                         `json:"user_id"`
		Email          string                      `json:"email"`
		DeliveryEmail  string                      `json:"delivery_email"`
		IsAdmin        bool                        `json:"is_admin"`
		IsGuest        bool                        `json:"is_guest"`
		IsBot          bool                        `json:"is_bot"`
		BotOwnerID     int                         `json:"bot_owner_id"`
		BotType        int                         `json:"bot_type"`
		IsActive       bool                        `json:"is_active"`
		IsOwner        bool                        `json:"is_owner"`
		IsBillingAdmin bool                        `json:"is_billing_admin"`
		Role           zulip.OrganizationRoleLevel `json:"role"`
		FullName       string                      `json:"full_name"`
		Timezone       string                      `json:"timezone"`
		DateJoined     string                      `json:"date_joined"`
		ProfileData    map[string]struct {
			Value         string `json:"value"`
			RenderedValue string `json:"rendered_value,omitempty"`
		} `json:"profile_data"`
		AvatarUrl     string `json:"avatar_url"`
		AvatarVersion int    `json:"avatar_version"`
	} `json:"members"`
}

func (aer *GetUsersResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &aer.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &aer.getUsersResponseData); err != nil {
		return err
	}

	return nil
}

type getUsersOptions struct {
	clientGravatar struct {
		fieldName string
		value     *bool
	}
	includeCustomProfileFields struct {
		fieldName string
		value     *bool
	}
}

type GetUsersOption func(*getUsersOptions)

func ClientGravatars(value bool) GetUsersOption {
	return func(args *getUsersOptions) {
		args.clientGravatar.fieldName = "client_gravatar"
		args.clientGravatar.value = &value
	}
}

func IncludeCustomProfilesFields(value bool) GetUsersOption {
	return func(args *getUsersOptions) {
		args.includeCustomProfileFields.fieldName = "include_custom_profile_fields"
		args.includeCustomProfileFields.value = &value
	}
}

func (svc *Service) GetUsers(ctx context.Context, options ...GetUsersOption) (*GetUsersResponse, error) {
	const (
		path   = "/api/v1/users"
		method = http.MethodGet
	)

	msg := map[string]any{}

	opts := getUsersOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.clientGravatar.value != nil {
		msg[opts.clientGravatar.fieldName] = *opts.clientGravatar.value
	}

	if opts.includeCustomProfileFields.value != nil {
		msg[opts.includeCustomProfileFields.fieldName] = *opts.includeCustomProfileFields.value
	}

	resp := GetUsersResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
