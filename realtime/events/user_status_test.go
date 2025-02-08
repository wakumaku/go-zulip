package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestUserStatus(t *testing.T) {
	eventExample := `{
    "away": true,
    "emoji_code": "1f697",
    "emoji_name": "car",
    "id": 0,
    "reaction_type": "unicode_emoji",
    "status_text": "out to lunch",
    "type": "user_status",
    "user_id": 10
}`

	v := events.UserStatus{}
	err := json.Unmarshal([]byte(eventExample), &v)
	assert.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.UserStatusType, v.EventType())

	assert.Equal(t, true, v.UserStatusData.Away)
	assert.Equal(t, "1f697", v.UserStatusData.EmojiCode)
	assert.Equal(t, "car", v.UserStatusData.EmojiName)
	assert.Equal(t, "unicode_emoji", v.UserStatusData.ReactionType)
	assert.Equal(t, "out to lunch", v.UserStatusData.StatusText)
	assert.Equal(t, 10, v.UserStatusData.UserID)
}
