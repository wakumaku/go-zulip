package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/messages"
)

func TestGetMessages(t *testing.T) {
	client := createMockClient(`{
    "anchor": 21,
    "found_anchor": true,
    "found_newest": true,
    "messages": [
        {
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
        {
            "avatar_url": "https://secure.gravatar.com/avatar/6d8cad0fd00256e7b40691d27ddfd466?d=identicon&version=1",
            "client": "populate_db",
            "content": "<p>Wait, is this from the frontend js code or backend python code</p>",
            "content_type": "text/html",
            "display_recipient": "Verona",
            "flags": [
                "read"
            ],
            "id": 21,
            "is_me_message": false,
            "reactions": [],
            "recipient_id": 20,
            "sender_email": "hamlet@zulip.com",
            "sender_full_name": "King Hamlet",
            "sender_id": 4,
            "sender_realm_str": "zulip",
            "stream_id": 5,
            "subject": "Verona3",
            "submessages": [],
            "timestamp": 1527939746,
            "topic_links": [],
            "type": "stream"
        }
    ],
    "msg": "",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)

	msg := map[string]any{
		"anchor":          "21",
		"include_anchor":  true,
		"num_before":      1,
		"num_after":       3,
		"narrow":          "[{\"operator\":\"channel\",\"operand\":\"Verona\"},{\"operator\":\"sender\",\"operand\":\"iago@zulip.com\"}]",
		"client_gravatar": true,
		"apply_markdown":  true,
		"message_ids":     []int{16, 21},
	}

	resp, err := messagesSvc.GetMessages(context.Background(),
		messages.Anchor("21"),
		messages.IncludeAnchor(true),
		messages.NumBefore(1),
		messages.NumAfter(3),
		messages.NarrowMessage(zulip.Narrower{}.
			Add(zulip.Channel, "Verona").
			Add(zulip.Sender, "iago@zulip.com")),
		messages.ClientGravatarMessage(true),
		messages.ApplyMarkdownMessage(true),
		messages.MessageIDs([]int{16, 21}),
	)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.IsSuccess())

	assert.Equal(t, 21, resp.Anchor)
	assert.Equal(t, true, resp.FoundAnchor)
	assert.Equal(t, true, resp.FoundNewest)
	assert.Equal(t, 2, len(resp.Messages))
	assert.Equal(t, "https://secure.gravatar.com/avatar/6d8cad0fd00256e7b40691d27ddfd466?d=identicon&version=1", resp.Messages[0].AvatarURL)
	assert.Equal(t, "populate_db", resp.Messages[0].Client)
	assert.Equal(t, "<p>Security experts agree that relational algorithms are an interesting new topic in the field of networking, and scholars concur.</p>", resp.Messages[0].Content)
	assert.Equal(t, "text/html", resp.Messages[0].ContentType)
	assert.Equal(t, 16, resp.Messages[0].ID)
	assert.Equal(t, false, resp.Messages[0].IsMeMessage)
	assert.Equal(t, 27, resp.Messages[0].RecipientID)
	assert.Equal(t, 1, len(resp.Messages[0].Flags))
	assert.Equal(t, "read", resp.Messages[0].Flags[0])
	assert.Equal(t, "hamlet@zulip.com", resp.Messages[0].DisplayRecipient.Users[0].Email)
	assert.Equal(t, "King Hamlet", resp.Messages[0].DisplayRecipient.Users[0].FullName)
	assert.Equal(t, 4, resp.Messages[0].DisplayRecipient.Users[0].ID)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
