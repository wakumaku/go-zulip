package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestUpdateMessage(t *testing.T) {
	eventExample := `{
    "content": "new content",
    "edit_timestamp": 1594825451,
    "flags": [],
    "id": 0,
    "is_me_message": false,
    "message_id": 58,
    "message_ids": [
        58,
        57
    ],
    "orig_content": "hello",
    "orig_rendered_content": "<p>hello</p>",
    "orig_subject": "test",
    "propagate_mode": "change_all",
    "rendered_content": "<p>new content</p>",
    "rendering_only": false,
    "stream_id": 5,
    "stream_name": "Verona",
    "subject": "new_topic",
    "topic_links": [],
    "type": "update_message",
    "user_id": 10
}`

	v := events.UpdateMessage{}
	err := json.Unmarshal([]byte(eventExample), &v)
	require.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.UpdateMessageType, v.EventType())
	assert.Equal(t, "update_message", v.EventOp())

	assert.Empty(t, v.Flags)
	assert.Equal(t, 58, v.MessageID)
	assert.Len(t, v.MessageIDs, 2)
	assert.Equal(t, "hello", *v.OrigContent)
	assert.Equal(t, "<p>hello</p>", *v.OrigRenderedContent)
	assert.Equal(t, "test", *v.OrigSubject)
	assert.Equal(t, "change_all", *v.PropagateMode)
	assert.Equal(t, "<p>new content</p>", *v.RenderedContent)
	assert.False(t, v.RenderingOnly)
	assert.Equal(t, 5, *v.StreamID)
	assert.Equal(t, "Verona", *v.StreamName)
	assert.Equal(t, "new_topic", *v.Subject)
	assert.Empty(t, v.TopicLinks)
	assert.Equal(t, "new content", *v.Content)
	assert.Equal(t, 1594825451, v.EditTimestamp)
	assert.False(t, *v.IsMeMessage)
	assert.Equal(t, 10, *v.UserID)
}
