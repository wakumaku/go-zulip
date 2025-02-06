package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestMessage(t *testing.T) {
	eventExample := `{
    "flags": [],
    "id": 1,
    "message": {
        "avatar_url": null,
        "client": "test suite",
        "content": "<p>First message ...<a href=\"user_uploads/2/ce/2Xpnnwgh8JWKxBXtTfD6BHKV/zulip.txt\">zulip.txt</a></p>",
        "content_type": "text/html",
        "display_recipient": "Denmark",
        "id": 31,
        "is_me_message": false,
        "reactions": [],
        "recipient_id": 23,
        "sender_email": "user10@zulip.testserver",
        "sender_full_name": "King Hamlet",
        "sender_id": 10,
        "sender_realm_str": "zulip",
        "stream_id": 1,
        "subject": "test",
        "submessages": [],
        "timestamp": 1594825416,
        "topic_links": [],
        "type": "stream"
    },
    "type": "message"
}`

	v := events.Message{}
	err := json.Unmarshal([]byte(eventExample), &v)
	assert.NoError(t, err)

	assert.Equal(t, 1, v.EventID())
	assert.Equal(t, events.MessageType, v.EventType())
	assert.Equal(t, "message", v.EventOp())

	assert.Equal(t, 0, len(v.Flags))
	assert.Equal(t, 31, v.Message.ID)
	assert.Equal(t, "test suite", v.Message.Client)
	assert.Equal(t, "<p>First message ...<a href=\"user_uploads/2/ce/2Xpnnwgh8JWKxBXtTfD6BHKV/zulip.txt\">zulip.txt</a></p>", v.Message.Content)
	assert.Equal(t, "text/html", v.Message.ContentType)
	assert.Equal(t, "Denmark", v.Message.DisplayRecipient.Channel)
	assert.Equal(t, 31, v.Message.ID)
	assert.Equal(t, false, v.Message.IsMeMessage)
	assert.Equal(t, 0, len(v.Message.Reactions))
	assert.Equal(t, 23, v.Message.RecipientID)
}
