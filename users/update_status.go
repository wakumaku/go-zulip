package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type UpdateStatusResponse struct {
	zulip.APIResponseBase
}

func (g *UpdateStatusResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	return nil
}

type updateStatusOptions struct {
	statusText struct {
		fieldName string
		value     *string
	}
	emojiName struct {
		fieldName string
		value     *string
	}
	emojiCode struct {
		fieldName string
		value     *string
	}
	reactionType struct {
		fieldName string
		value     *zulip.ReactionType
	}
}

type UpdateStatusOption func(*updateStatusOptions)

func StatusText(value string) UpdateStatusOption {
	return func(args *updateStatusOptions) {
		args.statusText.fieldName = "status_text"
		args.statusText.value = &value
	}
}

func StatusEmojiName(value string) UpdateStatusOption {
	return func(args *updateStatusOptions) {
		args.emojiName.fieldName = "emoji_name"
		args.emojiName.value = &value
	}
}

func StatusEmojiCode(value string) UpdateStatusOption {
	return func(args *updateStatusOptions) {
		args.emojiCode.fieldName = "emoji_code"
		args.emojiCode.value = &value
	}
}

func StatusReactionType(value zulip.ReactionType) UpdateStatusOption {
	return func(args *updateStatusOptions) {
		args.reactionType.fieldName = "reaction_type"
		args.reactionType.value = &value
	}
}

func (svc *Service) UpdateStatus(ctx context.Context, opts ...UpdateStatusOption) (*UpdateStatusResponse, error) {
	const (
		path   = "/api/v1/users/me/status"
		method = http.MethodPost
	)

	msg := map[string]any{}

	options := updateStatusOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.statusText.value != nil {
		msg[options.statusText.fieldName] = *options.statusText.value
	}

	if options.emojiName.value != nil {
		msg[options.emojiName.fieldName] = *options.emojiName.value
	}

	if options.emojiCode.value != nil {
		msg[options.emojiCode.fieldName] = *options.emojiCode.value
	}

	if options.reactionType.value != nil {
		msg[options.reactionType.fieldName] = *options.reactionType.value
	}

	resp := UpdateStatusResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
