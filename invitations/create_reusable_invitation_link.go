package invitations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wakumaku/go-zulip"
)

type CreateReusableInvitationLinkResponse struct {
	zulip.APIResponseBase
	createReusableInvitationLinkData
}

type createReusableInvitationLinkData struct {
	InviteLink string `json:"invite_link"`
}

func (c *CreateReusableInvitationLinkResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &c.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &c.createReusableInvitationLinkData); err != nil {
		return err
	}

	return nil
}

type createReusableInvitationLinkOptions struct {
	InviteExpiresInMinutes           *int
	InviteAs                         *zulip.OrganizationRoleLevel
	StreamIds                        *[]int
	IncludeRealmDefaultSubscriptions *bool
}

type CreateReusableInvitationLinkOption func(*createReusableInvitationLinkOptions)

func InviteExpiresInMinutes(minutes int) CreateReusableInvitationLinkOption {
	return func(o *createReusableInvitationLinkOptions) {
		o.InviteExpiresInMinutes = &minutes
	}
}

func InviteAs(roleLevel zulip.OrganizationRoleLevel) CreateReusableInvitationLinkOption {
	return func(o *createReusableInvitationLinkOptions) {
		o.InviteAs = &roleLevel
	}
}

func StreamIds(ids []int) CreateReusableInvitationLinkOption {
	return func(o *createReusableInvitationLinkOptions) {
		o.StreamIds = &ids
	}
}

func IncludeRealmDefaultSubscriptions(include bool) CreateReusableInvitationLinkOption {
	return func(o *createReusableInvitationLinkOptions) {
		o.IncludeRealmDefaultSubscriptions = &include
	}
}

func (svc *Service) CreateReusableInvitationLink(ctx context.Context, options ...CreateReusableInvitationLinkOption) (*CreateReusableInvitationLinkResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/invites/multiuse"
	)

	msg := map[string]any{}

	opts := createReusableInvitationLinkOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.InviteExpiresInMinutes != nil {
		msg["invite_expires_in_minutes"] = *opts.InviteExpiresInMinutes
	}

	if opts.InviteAs != nil {
		msg["invite_as"] = *opts.InviteAs
	}

	if opts.StreamIds != nil {
		msg["stream_ids"] = strings.ReplaceAll(fmt.Sprint(*opts.StreamIds), " ", ",")
	}

	if opts.IncludeRealmDefaultSubscriptions != nil {
		msg["include_realm_default_subscriptions"] = *opts.IncludeRealmDefaultSubscriptions
	}

	resp := CreateReusableInvitationLinkResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
