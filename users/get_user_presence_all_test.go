package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/users"
)

func TestGetUserPresenceAll(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "presences": {
        "iago@zulip.com": {
            "aggregated": {
                "client": "website",
                "status": "active",
                "timestamp": 1656958485
            },
            "website": {
                "client": "website",
                "pushable": false,
                "status": "active",
                "timestamp": 1656958485
            }
        }
    },
    "result": "success",
    "server_timestamp": 1656958539.6287155
}`)

	userSvc := users.NewService(client)

	resp, err := userSvc.GetUserPresenceAll(context.Background())
	require.NoError(t, err)

	iago := resp.Presences["iago@zulip.com"]
	assert.Equal(t, "active", iago.Aggregated.Status)
	assert.Equal(t, 1656958485, iago.Aggregated.Timestamp)
	assert.Equal(t, "active", iago.Website.Status)
	assert.Equal(t, 1656958485, iago.Website.Timestamp)
}
