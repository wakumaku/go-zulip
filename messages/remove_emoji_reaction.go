package messages

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/wakumaku/go-zulip"
)

type RemoveEmojiReactionResponse struct {
	zulip.APIResponseBase
}

type removeEmojiReactionOptions struct {
	emojiName    string
	emojiCode    string
	reactionType zulip.ReactionType
}

type RemoveEmojiReactionOption func(*removeEmojiReactionOptions) error

func RemoveEmojiReactionEmojiName(emojiName string) RemoveEmojiReactionOption {
	return func(o *removeEmojiReactionOptions) error {
		if strings.TrimSpace(emojiName) == "" {
			return errors.New("emoji name is empty")
		}
		o.emojiName = emojiName
		return nil
	}
}

func RemoveEmojiReactionEmojiCode(emojiCode string) RemoveEmojiReactionOption {
	return func(o *removeEmojiReactionOptions) error {
		if strings.TrimSpace(emojiCode) == "" {
			return errors.New("emoji code is empty")
		}
		o.emojiCode = emojiCode
		return nil
	}
}

func RemoveEmojiReactionReactionType(reactionType zulip.ReactionType) RemoveEmojiReactionOption {
	return func(o *removeEmojiReactionOptions) error {
		o.reactionType = reactionType
		return nil
	}
}

func (svc *Service) RemoveEmojiReaction(ctx context.Context, messageID int, options ...RemoveEmojiReactionOption) (*RemoveEmojiReactionResponse, error) {
	const (
		method = http.MethodDelete
		path   = "/api/v1/messages/{message_id}/reactions"
	)
	patchPath := strings.Replace(path, "{message_id}", fmt.Sprintf("%d", messageID), 1)

	msg := map[string]any{}

	opts := removeEmojiReactionOptions{}
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}

	if opts.emojiName != "" {
		msg["emoji_name"] = opts.emojiName
	}

	if opts.emojiCode != "" {
		msg["emoji_code"] = opts.emojiCode
	}

	if opts.reactionType != "" {
		msg["reaction_type"] = opts.reactionType
	}

	resp := RemoveEmojiReactionResponse{}
	if err := svc.client.DoRequest(ctx, method, patchPath, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
