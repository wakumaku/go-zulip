package channels_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/channels"
)

func TestGetStreamTopics(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success",
    "topics": [
        {"max_id": 100, "name": "topic A"},
        {"max_id": 98, "name": "topic B"}
    ]
}`)

	channelsSvc := channels.NewService(client)

	resp, err := channelsSvc.GetStreamTopics(context.Background(), 42)
	require.NoError(t, err)
	assert.True(t, resp.IsSuccess())

	assert.Equal(t, "/api/v1/users/me/42/topics", client.(*mockClient).path)
	assert.Equal(t, map[string]any{}, client.(*mockClient).paramsSent)

	require.Len(t, resp.Topics, 2)
	assert.Equal(t, "topic A", resp.Topics[0].Name)
	assert.Equal(t, 100, resp.Topics[0].MaxID)
}

func TestGetStreamTopicsAllowEmptyTopicName(t *testing.T) {
	client := createMockClient(`{"msg": "", "result": "success", "topics": []}`)
	channelsSvc := channels.NewService(client)

	_, err := channelsSvc.GetStreamTopics(context.Background(), 42, channels.AllowEmptyTopicName(true))
	require.NoError(t, err)

	assert.Equal(t, map[string]any{"allow_empty_topic_name": true}, client.(*mockClient).paramsSent)
}
