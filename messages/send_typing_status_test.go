package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/messages"
)

func TestSendTypingStatusDirect(t *testing.T) {
	client := createMockClient(`{"msg": "", "result": "success"}`)
	messagesSvc := messages.NewService(client)

	resp, err := messagesSvc.SendTypingStatus(context.Background(),
		messages.TypingStart,
		messages.ToDirect(10, 20),
	)
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

	assert.Equal(t, "/api/v1/typing", client.(*mockClient).path)
	assert.Equal(t, map[string]any{
		"op":   "start",
		"type": "direct",
		"to":   "[10,20]",
	}, client.(*mockClient).paramsSent)
}

func TestSendTypingStatusChannel(t *testing.T) {
	client := createMockClient(`{"msg": "", "result": "success"}`)
	messagesSvc := messages.NewService(client)

	resp, err := messagesSvc.SendTypingStatus(context.Background(),
		messages.TypingStop,
		messages.ToChannelTopic(42, "greetings"),
	)
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

	assert.Equal(t, "/api/v1/typing", client.(*mockClient).path)
	assert.Equal(t, map[string]any{
		"op":        "stop",
		"type":      "channel",
		"stream_id": 42,
		"topic":     "greetings",
	}, client.(*mockClient).paramsSent)
}

func TestSendTypingStatusDirectRequiresRecipients(t *testing.T) {
	client := createMockClient(`{"msg": "", "result": "success"}`)
	messagesSvc := messages.NewService(client)

	_, err := messagesSvc.SendTypingStatus(context.Background(), messages.TypingStart)
	require.Error(t, err)
}
