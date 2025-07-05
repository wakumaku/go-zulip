package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestRealmUser(t *testing.T) {
	eventExample := `{
    "id": 0,
    "op": "add",
    "person": {
        "avatar_url": "https://secure.gravatar.com/avatar/c6b5578d4964bd9c5fae593c6868912a?d=identicon&version=1",
        "avatar_version": 1,
        "date_joined": "2020-07-15T15:04:02.030833+00:00",
        "delivery_email": null,
        "email": "foo@zulip.com",
        "full_name": "full name",
        "is_active": true,
        "is_admin": false,
        "is_billing_admin": false,
        "is_bot": false,
        "is_guest": false,
        "is_owner": false,
        "profile_data": {},
        "role": 400,
        "timezone": "",
        "user_id": 38
    },
    "type": "realm_user"
}`

	v := events.RealmUser{}
	err := json.Unmarshal([]byte(eventExample), &v)
	require.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.RealmUserType, v.EventType())
	assert.Equal(t, "add", v.EventOp())

	assert.Equal(t, "https://secure.gravatar.com/avatar/c6b5578d4964bd9c5fae593c6868912a?d=identicon&version=1", v.Person.AvatarUrl)
	assert.Equal(t, 1, v.Person.AvatarVersion)
	assert.Equal(t, "2020-07-15T15:04:02.030833+00:00", v.Person.DateJoined)
	assert.Empty(t, v.Person.DeliveryEmail)
	assert.Equal(t, 38, v.Person.UserId)
}
