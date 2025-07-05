package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/users"
)

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name         string
		expectedPath string
		callFunc     func(userSvc *users.Service, profileData users.ProfileData) (*users.UpdateUserResponse, error)
	}{
		{
			name:         "UpdateUser by ID",
			expectedPath: "/api/v1/users/11",
			callFunc: func(userSvc *users.Service, profileData users.ProfileData) (*users.UpdateUserResponse, error) {
				return userSvc.UpdateUser(context.Background(), 11,
					users.FullName("King Hamlet"),
					users.Role(zulip.MemberRole),
					users.SetProfileData(profileData),
					users.NewEmail("newemail@xxx.com"),
				)
			},
		},
		{
			name:         "UpdateUser by Email",
			expectedPath: "/api/v1/users/test@tester.com",
			callFunc: func(userSvc *users.Service, profileData users.ProfileData) (*users.UpdateUserResponse, error) {
				return userSvc.UpdateUserByEmail(context.Background(), "test@tester.com",
					users.FullName("King Hamlet"),
					users.Role(zulip.MemberRole),
					users.SetProfileData(profileData),
					users.NewEmail("newemail@xxx.com"),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := createMockClient(`{
    "msg": "",
    "result": "success"
}`)

			userSvc := users.NewService(client)

			profileData := users.ProfileData{}
			profileData = append(profileData, users.ProfileDataItem{ID: 4, Value: "0"})
			profileData = append(profileData, users.ProfileDataItem{ID: 5, Value: "1909-04-05"})

			resp, err := tt.callFunc(userSvc, profileData)

			require.NoError(t, err)
			assert.Equal(t, "success", resp.Result())

			// validate the parameters sent are correct
			expectedMsg := map[string]interface{}{
				"full_name":    "King Hamlet",
				"role":         zulip.MemberRole,
				"profile_data": `[{"id":4,"value":"0"},{"id":5,"value":"1909-04-05"}]`,
				"new_email":    "newemail@xxx.com",
			}

			assert.Equal(t, tt.expectedPath, client.(*mockClient).path)
			assert.Equal(t, expectedMsg, client.(*mockClient).paramsSent)
		})
	}
}
