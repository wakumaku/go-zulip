package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/users"
)

func TestGetUser(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success",
    "user": {
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
    }
}`)

	userSvc := users.NewService(client)

	resp, err := userSvc.GetUser(context.Background(), 10,
		users.ClientGravatar(true),
		users.IncludeCustomProfileFields(true),
	)

	assert.NoError(t, err)
	assert.Equal(t, "https://secure.gravatar.com/avatar/6d8cad0fd00256e7b40691d27ddfd466?d=identicon&version=1", resp.User.AvatarUrl)
	assert.Equal(t, "King Hamlet", resp.User.FullName)
	assert.Equal(t, false, resp.User.IsAdmin)
	assert.Equal(t, 10, resp.User.UserID)
	assert.Equal(t, "hamlet@zulip.com", resp.User.Email)
	assert.Equal(t, zulip.MemberRole, resp.User.Role) // 400

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/users/10", client.(*mockClient).path)
}

func TestGetUserByEmail(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success",
    "user": {
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
    }
}`)

	userSvc := users.NewService(client)

	resp, err := userSvc.GetUserByEmail(context.Background(), "hamlet@zulip.com",
		users.ClientGravatar(true),
		users.IncludeCustomProfileFields(true),
	)
	assert.NoError(t, err)

	assert.Equal(t, "https://secure.gravatar.com/avatar/6d8cad0fd00256e7b40691d27ddfd466?d=identicon&version=1", resp.User.AvatarUrl)
	assert.Equal(t, "King Hamlet", resp.User.FullName)
	assert.Equal(t, false, resp.User.IsAdmin)
	assert.Equal(t, 10, resp.User.UserID)
	assert.Equal(t, "hamlet@zulip.com", resp.User.Email)
	assert.Equal(t, zulip.MemberRole, resp.User.Role) // 400

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/users/hamlet@zulip.com", client.(*mockClient).path)
}
