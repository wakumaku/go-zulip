package messages_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/messages"
)

func TestUpdatePersonalMessageFlags(t *testing.T) {
	client := createMockClient(`{
    "messages": [
        4,
        18,
        15
    ],
    "msg": "",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)
	resp, err := messagesSvc.UpdatePersonalMessageFlags(context.Background(), []int{4, 18, 15}, messages.OperationAdd, messages.FlagRead)
	require.NoError(t, err)

	assert.Equal(t, []int{4, 18, 15}, resp.Messages)

	// validate method & payload
	assert.Equal(t, http.MethodPost, client.(*mockClient).method)

	expectedParams := map[string]interface{}{
		"messages": "[4,18,15]",
		"op":       messages.OperationAdd,
		"flag":     messages.FlagRead,
	}
	assert.Equal(t, expectedParams, client.(*mockClient).paramsSent)
}
