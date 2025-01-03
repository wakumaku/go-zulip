package events

import "github.com/wakumaku/go-zulip"

const RealmUserType EventType = "realm_user"

type RealmUser struct {
	ID     int       `json:"id"`
	Op     string    `json:"op"`
	Type   EventType `json:"type"`
	Person Person    `json:"person"`
}

type Person struct {
	AvatarUrl      string                      `json:"avatar_url"`       //"avatar_url": "https://secure.gravatar.com/avatar/c6b5578d4964bd9c5fae593c6868912a?d=identicon&version=1",
	AvatarVersion  int                         `json:"avatar_version"`   //"avatar_version": 1,
	DateJoined     string                      `json:"date_joined"`      //"date_joined": "2020-07-15T15:04:02.030833+00:00",
	DeliveryEmail  string                      `json:"delivery_email"`   //"delivery_email": null,
	Email          string                      `json:"email"`            //"email": "foo@zulip.com",
	FullName       string                      `json:"full_name"`        //"full_name": "full name",
	IsActive       bool                        `json:"is_active"`        //"is_active": true,
	IsAdmin        bool                        `json:"is_admin"`         //"is_admin": false,
	IsBillingAdmin bool                        `json:"is_billing_admin"` //"is_billing_admin": false,
	IsBot          bool                        `json:"is_bot"`           //"is_bot": false,
	IsGuest        bool                        `json:"is_guest"`         //"is_guest": false,
	IsOwner        bool                        `json:"is_owner"`         //"is_owner": false,
	ProfileData    map[string]any              `json:"profile_data"`     //"profile_data": {},
	Role           zulip.OrganizationRoleLevel `json:"role"`             //"role": 400,
	Timezone       string                      `json:"timezone"`         //"timezone": "",
	UserId         int                         `json:"user_id"`          //"user_id": 38
}

func (e *RealmUser) EventID() int {
	return e.ID
}

func (e *RealmUser) EventType() EventType {
	return e.Type
}

func (e *RealmUser) EventOp() string {
	return e.Op
}
