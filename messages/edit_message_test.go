package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/messages"
)

func TestEditMessage(t *testing.T) {
	client := createMockClient(`{
    "detached_uploads": [
        {
            "create_time": 1687984706000,
            "id": 3,
            "messages": [],
            "name": "1253601-1.jpg",
            "path_id": "2/5d/BD5NRptFxPDKY3RUKwhhup8r/1253601-1.jpg",
            "size": 1339060
        }
    ],
    "msg": "",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)

	msg := map[string]any{
		"content":                         "new content",
		"propagate_mode":                  messages.PropagateModeAll,
		"send_notification_to_new_thread": true,
		"send_notification_to_old_thread": true,
		"stream_id":                       100,
		"topic":                           "new topic",
	}

	resp, err := messagesSvc.EditMessage(context.Background(),
		1,
		messages.MoveToTopic("new topic"),
		messages.SetPropagateMode(messages.PropagateModeAll),
		messages.SendNotificationToOldThread(true),
		messages.SendNotificationToNewThread(true),
		messages.NewContent("new content"),
		messages.SetStreamID(100),
	)
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

	assert.Len(t, resp.DetachedUploads, 1)
	assert.Equal(t, 1687984706000, resp.DetachedUploads[0].CreateTime)
	assert.Equal(t, 3, resp.DetachedUploads[0].ID)
	assert.Empty(t, resp.DetachedUploads[0].Messages)
	assert.Equal(t, "1253601-1.jpg", resp.DetachedUploads[0].Name)
	assert.Equal(t, "2/5d/BD5NRptFxPDKY3RUKwhhup8r/1253601-1.jpg", resp.DetachedUploads[0].PathID)
	assert.Equal(t, 1339060, resp.DetachedUploads[0].Size)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages/1", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
