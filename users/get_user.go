package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetUserResponse struct {
	zulip.APIResponseBase
	getUserResponseData
}

type getUserResponseData struct {
	User getUserMeResponseData `json:"user"`
}

func (g *GetUserResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getUserResponseData); err != nil {
		return err
	}

	return nil
}

type getUserOptions struct {
	clientGravatar struct {
		fieldName string
		value     *bool
	}
	includeCustomProfileFields struct {
		fieldName string
		value     *bool
	}
}

type GetUserOption func(*getUserOptions)

func ClientGravatar(value bool) GetUserOption {
	return func(args *getUserOptions) {
		args.clientGravatar.fieldName = "client_gravatar"
		args.clientGravatar.value = &value
	}
}

func IncludeCustomProfileFields(value bool) GetUserOption {
	return func(args *getUserOptions) {
		args.includeCustomProfileFields.fieldName = "include_custom_profile_fields"
		args.includeCustomProfileFields.value = &value
	}
}

func (svc *Service) GetUser(ctx context.Context, id int, options ...GetUserOption) (*GetUserResponse, error) {
	const (
		path = "/api/v1/users"
	)
	pathPatch := fmt.Sprintf("%s/%d", path, id)

	return svc.getUser(ctx, pathPatch, options...)
}

func (svc *Service) GetUserByEmail(ctx context.Context, email string, options ...GetUserOption) (*GetUserResponse, error) {
	const (
		path = "/api/v1/users"
	)
	pathPatch := fmt.Sprintf("%s/%s", path, email)

	return svc.getUser(ctx, pathPatch, options...)
}

func (svc *Service) getUser(ctx context.Context, path string, options ...GetUserOption) (*GetUserResponse, error) {
	const (
		method = http.MethodGet
	)

	msg := map[string]any{}

	opts := getUserOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.clientGravatar.value != nil {
		msg[opts.clientGravatar.fieldName] = *opts.clientGravatar.value
	}

	if opts.includeCustomProfileFields.value != nil {
		msg[opts.includeCustomProfileFields.fieldName] = *opts.includeCustomProfileFields.value
	}

	resp := GetUserResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
