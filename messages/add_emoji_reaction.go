package messages

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/wakumaku/go-zulip"
)

type AddEmojiReactionResponse struct {
	zulip.APIResponseBase
}

type addEmojiReactionOptions struct {
	emojiCode    string
	reactionType zulip.ReactionType
}

type AddEmojiReactionOption func(*addEmojiReactionOptions) error

func AddEmojiReactionEmojiCode(emojiCode string) AddEmojiReactionOption {
	return func(o *addEmojiReactionOptions) error {
		if strings.TrimSpace(emojiCode) == "" {
			return errors.New("emoji code is empty")
		}

		o.emojiCode = emojiCode

		return nil
	}
}

func AddEmojiReactionReactionType(reactionType zulip.ReactionType) AddEmojiReactionOption {
	return func(o *addEmojiReactionOptions) error {
		o.reactionType = reactionType
		return nil
	}
}

func (svc *Service) AddEmojiReaction(ctx context.Context, messageID int, emojiName string, options ...AddEmojiReactionOption) (*AddEmojiReactionResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/messages/{message_id}/reactions"
	)

	patchPath := strings.Replace(path, "{message_id}", fmt.Sprintf("%d", messageID), 1)

	msg := map[string]any{
		"emoji_name": emojiName,
	}

	opts := addEmojiReactionOptions{}
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}

	if opts.emojiCode != "" {
		msg["emoji_code"] = opts.emojiCode
	}

	if opts.reactionType != "" {
		msg["reaction_type"] = opts.reactionType
	}

	resp := AddEmojiReactionResponse{}
	if err := svc.client.DoRequest(ctx, method, patchPath, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
