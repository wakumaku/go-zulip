package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/users"
)

func TestUpdateStatus(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success"
}`)

	userSvc := users.NewService(client)

	resp, err := userSvc.UpdateStatus(context.Background(),
		users.StatusEmojiCode("1f697"),
		users.StatusEmojiName("car"),
		users.StatusReactionType("unicode_emoji"),
		users.StatusText("on vacation"),
	)
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Result())

	// validate the parameters sent are correct
	msg := map[string]interface{}{
		"emoji_code":    "1f697",
		"emoji_name":    "car",
		"reaction_type": zulip.UnicodeEmojiType,
		"status_text":   "on vacation",
	}
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
