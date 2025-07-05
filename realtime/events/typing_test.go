package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestTyping(t *testing.T) {
	eventExample := `{
    "id": 0,
    "message_type": "direct",
    "op": "start",
    "recipients": [
        {
            "email": "user8@zulip.testserver",
            "user_id": 8
        },
        {
            "email": "user10@zulip.testserver",
            "user_id": 10
        }
    ],
    "sender": {
        "email": "user10@zulip.testserver",
        "user_id": 10
    },
    "type": "typing"
}`

	v := events.Typing{}
	err := json.Unmarshal([]byte(eventExample), &v)
	require.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.TypingType, v.EventType())
	assert.Equal(t, "start", v.EventOp())

	assert.Equal(t, "direct", v.MessageType)

	assert.Equal(t, 8, v.Recipients[0].UserID)
	assert.Equal(t, "user8@zulip.testserver", v.Recipients[0].Email)

	assert.Equal(t, 10, v.Sender.UserID)
	assert.Equal(t, "user10@zulip.testserver", v.Sender.Email)
}
