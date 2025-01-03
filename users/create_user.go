package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type CreateUserResponse struct {
	zulip.APIResponseBase
	createUserResponseData
}

type createUserResponseData struct {
	UserID int `json:"user_id"`
}

func (aer *CreateUserResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &aer.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &aer.createUserResponseData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) CreateUser(ctx context.Context, email, password, fullName string) (*CreateUserResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/users"
	)

	msg := map[string]any{
		"email":     email,
		"password":  password,
		"full_name": fullName,
	}

	resp := CreateUserResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
