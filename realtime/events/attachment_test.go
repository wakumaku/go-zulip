package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestAttachment(t *testing.T) {
	eventExample := `{
    "attachment": {
        "create_time": 1594825414000,
        "id": 1,
        "messages": [],
        "name": "zulip.txt",
        "path_id": "2/ce/2Xpnnwgh8JWKxBXtTfD6BHKV/zulip.txt",
        "size": 6
    },
    "id": 0,
    "op": "add",
    "type": "attachment",
    "upload_space_used": 6
}`

	v := events.Attachment{}
	err := json.Unmarshal([]byte(eventExample), &v)
	require.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.AttachmentType, v.EventType())
	assert.Equal(t, "add", v.EventOp())

	assert.Equal(t, 1, v.Attachment.ID)
	assert.Equal(t, "zulip.txt", v.Attachment.Name)
	assert.Equal(t, "2/ce/2Xpnnwgh8JWKxBXtTfD6BHKV/zulip.txt", v.Attachment.PathID)
	assert.Equal(t, 6, v.Attachment.Size)
}
