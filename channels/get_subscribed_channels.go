package channels

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type GetSubscribedChannelsResponse struct {
	zulip.APIResponseBase
	getSubscribedChannelsResponseData
}

type getSubscribedChannelsResponseData struct {
	Subscriptions []struct {
		AudibleNotifications bool   `json:"audible_notifications"`
		Color                string `json:"color"`
		CreatorID            int    `json:"creator_id"`
		Description          string `json:"description"`
		DesktopNotifications bool   `json:"desktop_notifications"`
		InviteOnly           bool   `json:"invite_only"`
		IsArchived           bool   `json:"is_archived"`
		IsMuted              bool   `json:"is_muted"`
		Name                 string `json:"name"`
		PinToTop             bool   `json:"pin_to_top"`
		PushNotifications    bool   `json:"push_notifications"`
		StreamID             int    `json:"stream_id"`
		Subscribers          []int  `json:"subscribers"`
	} `json:"subscriptions"`
}

func (g *GetSubscribedChannelsResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getSubscribedChannelsResponseData); err != nil {
		return err
	}

	return nil
}

type getSubscribedChannelsOptions struct {
	includeSubscribers struct {
		fieldName string
		value     *bool
	}
}

type GetSubscribedChannelsOption func(*getSubscribedChannelsOptions)

func IncludeSubscribersList(includeSubscribers bool) GetSubscribedChannelsOption {
	return func(o *getSubscribedChannelsOptions) {
		o.includeSubscribers.fieldName = "include_subscribers"
		o.includeSubscribers.value = &includeSubscribers
	}
}

func (svc *Service) GetSubscribedChannels(ctx context.Context, options ...GetSubscribedChannelsOption) (*GetSubscribedChannelsResponse, error) {
	const (
		method = http.MethodGet
		path   = "/api/v1/users/me/subscriptions"
	)

	msg := map[string]any{}

	opts := getSubscribedChannelsOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.includeSubscribers.value != nil {
		msg[opts.includeSubscribers.fieldName] = *opts.includeSubscribers.value
	}

	resp := GetSubscribedChannelsResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
