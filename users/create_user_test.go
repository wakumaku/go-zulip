package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/users"
)

func TestCreateUser(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success",
    "user_id": 25
}`)

	userSvc := users.NewService(client)

	msg := map[string]any{
		"email":     "hello@world.test",
		"password":  "password",
		"full_name": "Hello World",
	}

	resp, err := userSvc.CreateUser(context.Background(), "hello@world.test", "password", "Hello World")
	require.NoError(t, err)
	assert.Equal(t, 25, resp.UserID)

	// validate the parameters sent are correct
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
