package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestUnknown(t *testing.T) {
	eventExample := `{
	"id": 0,
	"type": "new_event"
}`

	v := events.Unknown{}
	err := json.Unmarshal([]byte(eventExample), &v)
	require.NoError(t, err)

	assert.Equal(t, -1, v.EventID())
	assert.Equal(t, events.UnknownType, v.EventType())
	assert.Equal(t, "unknown", v.EventOp())
}
