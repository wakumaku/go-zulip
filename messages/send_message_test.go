package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/messages"
)

func TestSendMessageToChannel(t *testing.T) {
	client := createMockClient(`{
    "automatic_new_visibility_policy": 2,
    "id": 42,
    "msg": "",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)

	msg := map[string]any{
		"to":             messages.ToChannelName("channel"),
		"topic":          "topic",
		"type":           messages.ToChannel,
		"content":        "the message",
		"read_by_sender": true,
	}

	resp, err := messagesSvc.SendMessage(context.Background(),
		messages.ToChannelTopic(messages.ToChannelName("channel"), "topic"),
		"the message",
		messages.ToTopic("topic"), // repeated, it doesnt make sense on real usage when using ToChannelTopic
		messages.ReadBySender(true),
	)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.IsSuccess())

	assert.Equal(t, 42, resp.ID)
	assert.Equal(t, 2, resp.AutomaticNewVisibilityPolicy)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}

func TestSendMessageToUserIDs(t *testing.T) {
	client := createMockClient(`{
    "automatic_new_visibility_policy": 2,
    "id": 45,
    "msg": "",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)

	msg := map[string]any{
		"to":             "[1,2,3]",
		"type":           messages.ToDirect,
		"content":        "the message",
		"read_by_sender": true,
	}

	resp, err := messagesSvc.SendMessage(context.Background(),
		messages.ToUserIDs([]int{1, 2, 3}),
		"the message",
		messages.ReadBySender(true),
	)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.IsSuccess())

	assert.Equal(t, 45, resp.ID)
	assert.Equal(t, 2, resp.AutomaticNewVisibilityPolicy)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}

func TestSendMessageToUserName(t *testing.T) {
	client := createMockClient(`{
    "automatic_new_visibility_policy": 2,
    "id": 50,
    "msg": "",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)

	msg := map[string]any{
		"to":             messages.ToUserName("iago"),
		"type":           messages.ToDirect,
		"content":        "the message",
		"read_by_sender": true,
	}

	resp, err := messagesSvc.SendMessage(context.Background(),
		messages.ToUserName("iago"),
		"the message",
		messages.ReadBySender(true),
	)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.IsSuccess())

	assert.Equal(t, 50, resp.ID)
	assert.Equal(t, 2, resp.AutomaticNewVisibilityPolicy)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
