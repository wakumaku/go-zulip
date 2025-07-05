package specialty_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/specialty"
)

func TestFetchAPIKey(t *testing.T) {
	client := createMockClient(`{
    "api_key": "gjA04ZYcqXKalvYMA8OeXSfzUOLrtbZv",
    "email": "iago@zulip.com",
    "msg": "",
    "result": "success",
    "user_id": 5
}`)

	service := specialty.NewService(client)

	resp, err := service.FetchAPIKeyProduction(context.Background(), "iago@zulip.com", "password")
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Result())
	assert.Equal(t, "gjA04ZYcqXKalvYMA8OeXSfzUOLrtbZv", resp.APIKey)

	// validate the parameters sent are correct
	assert.Equal(t, map[string]any{
		"username": "iago@zulip.com",
		"password": "password",
	}, client.(*mockClient).paramsSent)
}

func TestFetchAPIKeyDevelopment(t *testing.T) {
	client := createMockClient(`{
    "api_key": "gjA04ZYcqXKalvYMA8OeXSfzUOLrtbZv",
    "email": "iago@zulip.com",
    "msg": "",
    "result": "success",
    "user_id": 5
}`)
	service := specialty.NewService(client)
	resp, err := service.FetchAPIKeyDevelopment(context.Background(), "iago@zulip.com")
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Result())
	assert.Equal(t, "gjA04ZYcqXKalvYMA8OeXSfzUOLrtbZv", resp.APIKey)

	// validate the parameters sent are correct
	assert.Equal(t, map[string]any{
		"username": "iago@zulip.com",
	}, client.(*mockClient).paramsSent)
}
