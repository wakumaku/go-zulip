package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type UpdateUserResponse struct {
	zulip.APIResponseBase
	updateUserResponseData
}

type updateUserResponseData struct{}

func (u *UpdateUserResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &u.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &u.updateUserResponseData); err != nil {
		return err
	}

	return nil
}

type updateUserOptions struct {
	fullName struct {
		fieldName string
		value     *string
	}
	role struct {
		fieldName string
		value     *zulip.OrganizationRoleLevel
	}
	profileData struct {
		fieldName string
		value     *ProfileData
	}
	newEmail struct {
		fieldName string
		value     *string
	}
}

type UpdateUserOption func(*updateUserOptions)

func FullName(value string) UpdateUserOption {
	return func(args *updateUserOptions) {
		args.fullName.fieldName = "full_name"
		args.fullName.value = &value
	}
}

func Role(value zulip.OrganizationRoleLevel) UpdateUserOption {
	return func(args *updateUserOptions) {
		args.role.fieldName = "role"
		args.role.value = &value
	}
}

type ProfileData []ProfileDataItem

type ProfileDataItem struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

func SetProfileData(value ProfileData) UpdateUserOption {
	return func(args *updateUserOptions) {
		args.profileData.fieldName = "profile_data"
		args.profileData.value = &value
	}
}

func NewEmail(value string) UpdateUserOption {
	return func(args *updateUserOptions) {
		args.newEmail.fieldName = "new_email"
		args.newEmail.value = &value
	}
}

func (svc *Service) UpdateUser(ctx context.Context, id int, options ...UpdateUserOption) (*UpdateUserResponse, error) {
	const (
		path = "/api/v1/users"
	)

	patchPath := fmt.Sprintf("%s/%d", path, id)

	return svc.updateUser(ctx, patchPath, options...)
}

func (svc *Service) UpdateUserByEmail(ctx context.Context, email string, options ...UpdateUserOption) (*UpdateUserResponse, error) {
	const (
		path = "/api/v1/users"
	)

	patchPath := fmt.Sprintf("%s/%s", path, email)

	return svc.updateUser(ctx, patchPath, options...)
}

func (svc *Service) updateUser(ctx context.Context, patchPath string, options ...UpdateUserOption) (*UpdateUserResponse, error) {
	const (
		method = http.MethodPatch
	)

	msg := map[string]any{}

	opts := updateUserOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.fullName.value != nil {
		msg[opts.fullName.fieldName] = *opts.fullName.value
	}

	if opts.role.value != nil {
		msg[opts.role.fieldName] = *opts.role.value
	}

	if opts.profileData.value != nil {
		jsonData, err := json.Marshal(opts.profileData.value)
		if err != nil {
			return nil, err
		}

		msg[opts.profileData.fieldName] = string(jsonData)
	}

	if opts.newEmail.value != nil {
		msg[opts.newEmail.fieldName] = *opts.newEmail.value
	}

	resp := UpdateUserResponse{}
	if err := svc.client.DoRequest(ctx, method, patchPath, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
