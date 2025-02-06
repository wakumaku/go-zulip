package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestSubmessage(t *testing.T) {
	eventExample := `{
    "content": "{\"type\":\"vote\",\"key\":\"58,1\",\"vote\":1}",
    "id": 28,
    "message_id": 970461,
    "msg_type": "widget",
    "sender_id": 58,
    "submessage_id": 4737,
    "type": "submessage"
}`

	v := events.Submessage{}
	err := json.Unmarshal([]byte(eventExample), &v)
	assert.NoError(t, err)

	assert.Equal(t, 28, v.EventID())
	assert.Equal(t, events.SubmessageType, v.EventType())
	assert.Equal(t, "submessage", v.EventOp())

	assert.Equal(t, "{\"type\":\"vote\",\"key\":\"58,1\",\"vote\":1}", v.Content)
	assert.Equal(t, 970461, v.MessageID)
	assert.Equal(t, "widget", v.MsgType)
	assert.Equal(t, 58, v.SenderID)
	assert.Equal(t, 4737, v.SubmessageID)
}
