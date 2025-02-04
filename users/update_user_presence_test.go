package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/users"
)

func TestUpdateUserPresence(t *testing.T) {
	client := createMockClient(`{
		"msg": "",
		"presence_last_update_id": 1000,
		"presences": {
			"10": {
				"active_timestamp": 1656958520,
				"idle_timestamp": 1656958530
			}
		},
		"result": "success",
		"server_timestamp": 1656958539.6287155
	}`)

	userSvc := users.NewService(client)

	resp, err := userSvc.UpdateUserPresence(context.Background(), users.UserPresenceActive)
	assert.NoError(t, err)
	assert.Equal(t, 1000, resp.PresenceLastUpdateID)
	assert.Equal(t, 1656958520, resp.Presences["10"].ActiveTimestamp)
	assert.Equal(t, 1656958530, resp.Presences["10"].IdleTimestamp)
	assert.Equal(t, 1656958539.6287155, resp.ServerTimestamp)
}
