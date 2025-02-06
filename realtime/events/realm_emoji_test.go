package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestRealmEmoji(t *testing.T) {
	eventExample := `{
    "id": 0,
    "op": "update",
    "realm_emoji": {
        "1": {
            "author_id": 11,
            "deactivated": false,
            "id": "1",
            "name": "green_tick",
            "source_url": "/user_avatars/2/emoji/images/1.png"
        },
        "2": {
            "author_id": 11,
            "deactivated": true,
            "id": "2",
            "name": "my_emoji",
            "source_url": "/user_avatars/2/emoji/images/2.png"
        }
    },
    "type": "realm_emoji"
}`

	v := events.RealmEmoji{}
	err := json.Unmarshal([]byte(eventExample), &v)
	assert.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.RealmEmojiType, v.EventType())
	assert.Equal(t, "update", v.EventOp())

	assert.Equal(t, 11, v.RealmEmoji["1"].AuthorID)
	assert.Equal(t, false, v.RealmEmoji["1"].Deactivated)
	assert.Equal(t, "1", v.RealmEmoji["1"].ID)
	assert.Equal(t, "green_tick", v.RealmEmoji["1"].Name)
	assert.Equal(t, "/user_avatars/2/emoji/images/1.png", v.RealmEmoji["1"].SourceURL)
}
