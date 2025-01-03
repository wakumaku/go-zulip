package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/messages"
)

func TestRemoveEmojiReaction(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)

	msg := map[string]any{
		"emoji_name":    "smile",
		"emoji_code":    "1f604",
		"reaction_type": zulip.ReactionUnicodeEmoji,
	}

	resp, err := messagesSvc.RemoveEmojiReaction(context.Background(),
		105,
		messages.RemoveEmojiReactionEmojiName("smile"),
		messages.RemoveEmojiReactionEmojiCode("1f604"),
		messages.RemoveEmojiReactionReactionType(zulip.ReactionUnicodeEmoji),
	)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.IsSuccess())

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages/105/reactions", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
