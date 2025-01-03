package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/messages"
)

func TestDeleteMessage(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success"
}`)

	messagesSvc := messages.NewService(client)

	resp, err := messagesSvc.DeleteMessage(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.IsSuccess())

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/messages/1", client.(*mockClient).path)
}
