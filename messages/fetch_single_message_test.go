package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/messages"
)

func TestFetchSingleMessage(t *testing.T) {
	client := createMockClient(`{
    "message": {
        "avatar_url": "https://secure.gravatar.com/avatar/6d8cad0fd00256e7b40691d27ddfd466?d=identicon&version=1",
        "client": "populate_db",
        "content": "<p>Security experts agree that relational algorithms are an interesting new topic in the field of networking, and scholars concur.</p>",
        "content_type": "text/html",
        "display_recipient": [
            {
                "email": "hamlet@zulip.com",
                "full_name": "King Hamlet",
                "id": 4,
                "is_mirror_dummy": false
            },
            {
                "email": "iago@zulip.com",
                "full_name": "Iago",
                "id": 5,
                "is_mirror_dummy": false
            },
            {
                "email": "prospero@zulip.com",
                "full_name": "Prospero from The Tempest",
                "id": 8,
                "is_mirror_dummy": false
            }
        ],
        "flags": [
            "read"
        ],
        "id": 16,
        "is_me_message": false,
        "reactions": [],
        "recipient_id": 27,
        "sender_email": "hamlet@zulip.com",
        "sender_full_name": "King Hamlet",
        "sender_id": 4,
        "sender_realm_str": "zulip",
        "subject": "",
        "submessages": [],
        "timestamp": 1527921326,
        "topic_links": [],
        "type": "private"
    },
    "msg": "",
    "raw_content": "**Don't** forget your towel!",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)

	msg := map[string]any{
		"apply_markdown": true,
	}

	resp, err := messagesSvc.FetchSingleMessage(context.Background(),
		16,
		messages.ApplyMarkdownSingleMessage(true),
	)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.IsSuccess())

	assert.Equal(t, "https://secure.gravatar.com/avatar/6d8cad0fd00256e7b40691d27ddfd466?d=identicon&version=1", resp.Message.AvatarURL)
	assert.Equal(t, "populate_db", resp.Message.Client)
	assert.Equal(t, "<p>Security experts agree that relational algorithms are an interesting new topic in the field of networking, and scholars concur.</p>", resp.Message.Content)
	assert.Equal(t, "text/html", resp.Message.ContentType)
	assert.Equal(t, 16, resp.Message.ID)
	assert.Equal(t, false, resp.Message.IsMeMessage)
	assert.Equal(t, 27, resp.Message.RecipientID)
	assert.Equal(t, "**Don't** forget your towel!", resp.RawContent)
	assert.Equal(t, 1, len(resp.Message.Flags))
	assert.Equal(t, "read", resp.Message.Flags[0])
	assert.Equal(t, false, resp.Message.DisplayRecipient.IsChannel)
	assert.Equal(t, 3, len(resp.Message.DisplayRecipient.Users))
	assert.Equal(t, "hamlet@zulip.com", resp.Message.DisplayRecipient.Users[0].Email)
	assert.Equal(t, "King Hamlet", resp.Message.DisplayRecipient.Users[0].FullName)
	assert.Equal(t, 4, resp.Message.DisplayRecipient.Users[0].ID)
	assert.Equal(t, false, resp.Message.DisplayRecipient.Users[0].IsMirrorDummy)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages/16", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
