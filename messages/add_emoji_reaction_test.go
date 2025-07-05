package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/messages"
)

func TestAddEmojiReaction(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)

	msg := map[string]any{
		"emoji_name":    "smile",
		"emoji_code":    "1f604",
		"reaction_type": zulip.UnicodeEmojiType,
	}

	resp, err := messagesSvc.AddEmojiReaction(context.Background(),
		25,
		"smile",
		messages.AddEmojiReactionEmojiCode("1f604"),
		messages.AddEmojiReactionReactionType(zulip.UnicodeEmojiType),
	)
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

	// validate the parameters sent are correct
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
