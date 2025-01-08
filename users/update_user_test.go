package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/users"
)

func TestUpdateUser(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success"
}`)

	userSvc := users.NewService(client)

	profileData := users.ProfileData{}
	profileData = append(profileData, users.ProfileDataItem{ID: 4, Value: "0"})
	profileData = append(profileData, users.ProfileDataItem{ID: 5, Value: "1909-04-05"})

	resp, err := userSvc.UpdateUser(context.Background(), 11,
		users.FullName("King Hamlet"),
		users.Role(zulip.MemberRole),
		users.SetProfileData(profileData),
		users.NewEmail("newemail@xxx.com"),
	)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Result())

	// validate the parameters sent are correct
	msg := map[string]interface{}{
		"full_name":    "King Hamlet",
		"role":         zulip.MemberRole,
		"profile_data": `[{"id":4,"value":"0"},{"id":5,"value":"1909-04-05"}]`,
		"new_email":    "newemail@xxx.com",
	}
	assert.Equal(t, "/api/v1/users/11", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
