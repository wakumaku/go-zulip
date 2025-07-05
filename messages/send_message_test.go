package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/messages"
	"github.com/wakumaku/go-zulip/messages/recipient"
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
		"to":             recipient.ToChannel("channel"),
		"topic":          "topic",
		"type":           "channel",
		"content":        "the message",
		"read_by_sender": true,
	}

	// Send message to a channel with the topic as parameter
	resp, err := messagesSvc.SendMessage(context.Background(),
		recipient.ToChannel("channel"),
		"the message",
		messages.ToTopic("topic"), // repeated, it doesnt make sense on real usage when using ToChannelTopic
		messages.ReadBySender(true),
	)
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

	assert.Equal(t, 42, resp.ID)
	assert.Equal(t, 2, resp.AutomaticNewVisibilityPolicy)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)

	// Using SendMessageToChannelTopic
	resp, err = messagesSvc.SendMessageToChannelTopic(context.Background(),
		recipient.ToChannel("channel"), "topic",
		"the message",
		messages.ReadBySender(true))
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

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
		"type":           "direct",
		"content":        "the message",
		"read_by_sender": true,
	}

	resp, err := messagesSvc.SendMessage(context.Background(),
		recipient.ToUsers([]int{1, 2, 3}),
		"the message",
		messages.ReadBySender(true),
	)
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

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
		"to":             "[\"iago\"]",
		"type":           "direct",
		"content":        "the message",
		"read_by_sender": true,
	}

	resp, err := messagesSvc.SendMessageToUsers(context.Background(),
		recipient.ToUser("iago"),
		"the message",
		messages.ReadBySender(true),
	)
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

	assert.Equal(t, 50, resp.ID)
	assert.Equal(t, 2, resp.AutomaticNewVisibilityPolicy)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}

func TestSendMessageToUserNames(t *testing.T) {
	client := createMockClient(`{
	"automatic_new_visibility_policy": 2,
	"id": 55,
	"msg": "",
	"result": "success"
}`)

	messagesSvc := messages.NewService(client)

	msg := map[string]any{
		"to":             "[\"iago\",\"cordelia\"]",
		"type":           "direct",
		"content":        "the message",
		"read_by_sender": true,
	}

	resp, err := messagesSvc.SendMessageToUsers(context.Background(),
		recipient.ToUsers([]string{"iago", "cordelia"}),
		"the message",
		messages.ReadBySender(true),
	)
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

	assert.Equal(t, 55, resp.ID)
	assert.Equal(t, 2, resp.AutomaticNewVisibilityPolicy)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
