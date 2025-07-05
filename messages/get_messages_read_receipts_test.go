package messages_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/messages"
)

func TestGetMessagesReadReceipts(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success",
    "user_ids": [
        3,
        7,
        9
    ]
}`)

	messagesSvc := messages.NewService(client)
	resp, err := messagesSvc.GetMessagesReadReceipts(context.Background(), 42)
	require.NoError(t, err)
	assert.Equal(t, []int{3, 7, 9}, resp.UserIDs)

	// validate method & path
	assert.Equal(t, http.MethodGet, client.(*mockClient).method)
	assert.Equal(t, "/api/v1/messages/42/read_receipts", client.(*mockClient).path)
}
