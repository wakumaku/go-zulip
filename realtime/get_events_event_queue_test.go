package realtime

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestGetEventsEventQueueResponse(t *testing.T) {
	data, err := os.ReadFile("testdata/events.json")
	assert.NoError(t, err)

	geeqr := GetEventsEventQueueResponse{}
	err = geeqr.UnmarshalJSON(data)
	assert.NoError(t, err)

	assert.Len(t, geeqr.Events, 13)

	assert.IsType(t, &events.Message{}, geeqr.Events[0])

	message := geeqr.Events[0].(*events.Message)
	assert.Equal(t, events.MessageType, message.EventType())
	assert.Equal(t, "go-zulip", message.Message.Client)
	assert.Equal(t, "greetings", message.Message.Subject)
	assert.Equal(t, "stream", message.Message.Type)

	attachment := geeqr.Events[1].(*events.Attachment)
	assert.Equal(t, events.AttachmentType, attachment.EventType())
	assert.Equal(t, "add", attachment.Op)
	assert.Equal(t, 1594825414000, attachment.Attachment.CreateTime)
	assert.Equal(t, "2/ce/2Xpnnwgh8JWKxBXtTfD6BHKV/zulip.txt", attachment.Attachment.PathID)

	attachment2 := geeqr.Events[2].(*events.Attachment)
	assert.Equal(t, events.AttachmentType, attachment2.EventType())
	assert.Equal(t, "update", attachment2.Op)

	attachment3 := geeqr.Events[3].(*events.Attachment)
	assert.Equal(t, events.AttachmentType, attachment3.EventType())
	assert.Equal(t, "remove", attachment3.Op)

	heartbeat := geeqr.Events[4].(*events.Heartbeat)
	assert.Equal(t, events.HeartbeatType, heartbeat.EventType())

	presence := geeqr.Events[5].(*events.Presence)
	assert.Equal(t, events.PresenceType, presence.EventType())
	assert.Equal(t, "user10@zulip.testserver", presence.Email)
	assert.Equal(t, 10, presence.UserID)
	assert.Equal(t, "idle", presence.Presence.Website.Status)

	realmEmoji := geeqr.Events[6].(*events.RealmEmoji)
	assert.Equal(t, events.RealmEmojiType, realmEmoji.EventType())
	assert.Equal(t, "update", realmEmoji.Op)
	assert.Equal(t, "green_tick", realmEmoji.RealmEmoji["1"].Name)
	assert.Equal(t, "/user_avatars/2/emoji/images/2.png", realmEmoji.RealmEmoji["2"].SourceURL)
}
