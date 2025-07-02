package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestDeleteMessage(t *testing.T) {
	eventExample := `{
    "id": 0,
	"type": "delete_message",
    "message_id": 58,
    "message_ids": [
        58,
        57
    ],
    "stream_id": 5,
    "topic": "new_topic"
}`

	v := events.DeleteMessage{}
	err := json.Unmarshal([]byte(eventExample), &v)
	assert.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.DeleteMessageType, v.EventType())
	assert.Equal(t, "delete_message", v.EventOp())

	assert.Equal(t, 58, *v.MessageID)
	assert.Equal(t, 2, len(v.MessageIDs))
	assert.Equal(t, 5, *v.StreamID)
	assert.Equal(t, "new_topic", *v.Topic)
}
