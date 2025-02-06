package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestHeartbeat(t *testing.T) {
	eventExample := `{
    "id": 0,
    "type": "heartbeat"
}`

	v := events.Heartbeat{}
	err := json.Unmarshal([]byte(eventExample), &v)
	assert.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.HeartbeatType, v.EventType())
	assert.Equal(t, "heartbeat", v.EventOp())
}
