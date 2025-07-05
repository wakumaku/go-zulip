package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/users"
)

func TestGetUsers(t *testing.T) {
	client := createMockClient(`{
    "members": [
        {
            "avatar_url": "https://secure.gravatar.com/avatar/818c212b9f8830dfef491b3f7da99a14?d=identicon&version=1",
            "bot_type": null,
            "date_joined": "2019-10-20T07:50:53.728864+00:00",
            "delivery_email": null,
            "email": "AARON@zulip.com",
            "full_name": "aaron",
            "is_active": true,
            "is_admin": false,
            "is_billing_admin": false,
            "is_bot": false,
            "is_guest": false,
            "is_owner": false,
            "profile_data": {},
            "role": 400,
            "timezone": "",
            "user_id": 7
        },
        {
            "avatar_url": "https://secure.gravatar.com/avatar/6d8cad0fd00256e7b40691d27ddfd466?d=identicon&version=1",
            "bot_type": null,
            "date_joined": "2019-10-20T07:50:53.729659+00:00",
            "delivery_email": null,
            "email": "hamlet@zulip.com",
            "full_name": "King Hamlet",
            "is_active": true,
            "is_admin": false,
            "is_billing_admin": false,
            "is_bot": false,
            "is_guest": false,
            "is_owner": false,
            "profile_data": {
                "1": {
                    "rendered_value": "<p>+0-11-23-456-7890</p>",
                    "value": "+0-11-23-456-7890"
                },
                "2": {
                    "rendered_value": "<p>I am:</p>\n<ul>\n<li>The prince of Denmark</li>\n<li>Nephew to the usurping Claudius</li>\n</ul>",
                    "value": "I am:\n* The prince of Denmark\n* Nephew to the usurping Claudius"
                },
                "3": {
                    "rendered_value": "<p>Dark chocolate</p>",
                    "value": "Dark chocolate"
                },
                "4": {
                    "value": "0"
                },
                "5": {
                    "value": "1900-01-01"
                },
                "6": {
                    "value": "https://blog.zulig.org"
                },
                "7": {
                    "value": "[11]"
                },
                "8": {
                    "value": "zulipbot"
                }
            },
            "role": 400,
            "timezone": "",
            "user_id": 10
        },
        {
            "avatar_url": "https://secure.gravatar.com/avatar/7328586831cdbb1627649bd857b1ee8c?d=identicon&version=1",
            "bot_owner_id": 11,
            "bot_type": 1,
            "date_joined": "2019-10-20T12:52:17.862053+00:00",
            "delivery_email": "iago-bot@zulipdev.com",
            "email": "iago-bot@zulipdev.com",
            "full_name": "Iago's Bot",
            "is_active": true,
            "is_admin": false,
            "is_billing_admin": false,
            "is_bot": true,
            "is_guest": false,
            "is_owner": false,
            "role": 400,
            "timezone": "",
            "user_id": 23
        }
    ],
    "msg": "",
    "result": "success"
}`)

	userSvc := users.NewService(client)

	resp, err := userSvc.GetUsers(context.Background(),
		users.ClientGravatars(true),
		users.IncludeCustomProfilesFields(true),
	)
	require.NoError(t, err)

	assert.Len(t, resp.Members, 3)

	assert.Equal(t, "https://secure.gravatar.com/avatar/818c212b9f8830dfef491b3f7da99a14?d=identicon&version=1", resp.Members[0].AvatarUrl)
	assert.Equal(t, "aaron", resp.Members[0].FullName)
	assert.False(t, resp.Members[0].IsAdmin)
	assert.Equal(t, 7, resp.Members[0].UserID)

	assert.Equal(t, 11, resp.Members[2].BotOwnerID)
	assert.Equal(t, 1, resp.Members[2].BotType)

	// validate the parameters sent are correct
	msg := map[string]interface{}{
		"client_gravatar":               true,
		"include_custom_profile_fields": true,
	}
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
	assert.Equal(t, "/api/v1/users", client.(*mockClient).path)
}
