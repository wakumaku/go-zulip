package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestPresence(t *testing.T) {
	eventExample := `{
    "email": "user10@zulip.testserver",
    "id": 0,
    "presence": {
        "website": {
            "client": "website",
            "pushable": false,
            "status": "idle",
            "timestamp": 1594825445
        }
    },
    "server_timestamp": 1594825445.3200784,
    "type": "presence",
    "user_id": 10
}`

	v := events.Presence{}
	err := json.Unmarshal([]byte(eventExample), &v)
	assert.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.PresenceType, v.EventType())
	assert.Equal(t, "presence", v.EventOp())

	assert.Equal(t, "idle", v.Presence.Website.Status)
}
