package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/users"
)

func TestGetUserPresence(t *testing.T) {
	client := createMockClient(`{
	"msg": "",
	"result": "success",
	"presence": {
		"aggregated": {
			"status": "active",
			"timestamp": 1590000000
		},
		"website": {
			"status": "active",
			"timestamp": 1590000000
		}
	}
}`)

	userSvc := users.NewService(client)

	resp, err := userSvc.GetUserPresence(context.Background(), "test@example.com")
	require.NoError(t, err)
	assert.Equal(t, "active", resp.Presence.Aggregated.Status)
	assert.Equal(t, 1590000000, resp.Presence.Aggregated.Timestamp)
	assert.Equal(t, "active", resp.Presence.Website.Status)
	assert.Equal(t, 1590000000, resp.Presence.Website.Timestamp)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/users/test@example.com/presence", client.(*mockClient).path)
}
