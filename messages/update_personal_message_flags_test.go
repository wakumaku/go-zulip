package messages_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/messages"
)

func TestUpdatePersonalMessageFlags(T *testing.T) {
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
	assert.NoError(T, err)

	assert.Equal(T, resp.Messages, []int{4, 18, 15})

	// validate method & payload
	assert.Equal(T, http.MethodPost, client.(*mockClient).method)
	expedtedParams := map[string]interface{}{
		"messages": "[4,18,15]",
		"op":       messages.OperationAdd,
		"flag":     messages.FlagRead,
	}
	assert.Equal(T, expedtedParams, client.(*mockClient).paramsSent)
}
