package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/users"
)

func TestGetUserStatus(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success",
    "status": {
        "emoji_code": "1f697",
        "emoji_name": "car",
        "reaction_type": "unicode_emoji",
        "status_text": "on vacation"
    }
}`)

	userSvc := users.NewService(client)

	resp, err := userSvc.GetUserStatus(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "1f697", resp.Status.EmojiCode)
	assert.Equal(t, "car", resp.Status.EmojiName)
	assert.Equal(t, zulip.ReactionType("unicode_emoji"), resp.Status.ReactionType)
	assert.Equal(t, "on vacation", resp.Status.StatusText)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/users/1/status", client.(*mockClient).path)
}
