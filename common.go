package zulip

type ReactionType string

const (
	UnicodeEmojiType    ReactionType = "unicode_emoji"
	RealmEmojiType      ReactionType = "realm_emoji"
	ZulipExtraEmojiType ReactionType = "zulip_extra_emoji"
)

type OrganizationRoleLevel int

const (
	OwnerRole         OrganizationRoleLevel = 100
	AdministratorRole OrganizationRoleLevel = 200
	ModeratorRole     OrganizationRoleLevel = 300
	MemberRole        OrganizationRoleLevel = 400
	GuestRole         OrganizationRoleLevel = 600
)
