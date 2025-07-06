package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/users"
)

func TestGetUserMe(t *testing.T) {
	client := createMockClient(`{
    "avatar_url": "https://secure.gravatar.com/avatar/af4f06322c177ef4e1e9b2c424986b54?d=identicon&version=1",
    "avatar_version": 1,
    "date_joined": "2019-10-20T07:50:53.728864+00:00",
    "delivery_email": "iago@zulip.com",
    "email": "iago@zulip.com",
    "full_name": "Iago",
    "is_active": true,
    "is_admin": true,
    "is_billing_admin": false,
    "is_bot": false,
    "is_guest": false,
    "is_owner": false,
    "max_message_id": 30,
    "msg": "",
    "profile_data": {
        "1": {
            "rendered_value": "<p>+1-234-567-8901</p>",
            "value": "+1-234-567-8901"
        },
        "2": {
            "rendered_value": "<p>Betrayer of Othello.</p>",
            "value": "Betrayer of Othello."
        },
        "3": {
            "rendered_value": "<p>Apples</p>",
            "value": "Apples"
        },
        "4": {
            "value": "emacs"
        },
        "5": {
            "value": "2000-01-01"
        },
        "6": {
            "value": "https://zulip.readthedocs.io/en/latest/"
        },
        "7": {
            "value": "[10]"
        },
        "8": {
            "value": "zulip"
        }
    },
    "result": "success",
    "role": 200,
    "timezone": "",
    "user_id": 5
}`)

	userSvc := users.NewService(client)

	resp, err := userSvc.GetUserMe(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 5, resp.UserID)
	assert.Equal(t, "Iago", resp.FullName)
	assert.True(t, resp.IsAdmin)
	assert.False(t, resp.IsBillingAdmin)
	assert.False(t, resp.IsBot)
	assert.False(t, resp.IsGuest)
	assert.False(t, resp.IsOwner)
	assert.Equal(t, 30, resp.MaxMessageID)
	assert.Equal(t, "https://secure.gravatar.com/avatar/af4f06322c177ef4e1e9b2c424986b54?d=identicon&version=1", resp.AvatarURL)
	assert.Equal(t, 1, resp.AvatarVersion)
	assert.Equal(t, "2019-10-20T07:50:53.728864+00:00", resp.DateJoined)
	assert.Equal(t, "iago@zulip.com", resp.DeliveryEmail)
	assert.Equal(t, "iago@zulip.com", resp.Email)
	assert.Equal(t, "Apples", resp.ProfileData["3"].Value)
}
