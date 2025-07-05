package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/messages"
)

func TestRenderAMessage(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "rendered": "<p><strong>foo</strong></p>",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)

	msg := map[string]any{
		"content": "**foo**",
	}

	resp, err := messagesSvc.RenderAMessage(context.Background(), "**foo**")
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages/render", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
