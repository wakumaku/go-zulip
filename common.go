package zulip

type ReactionType string

const (
	ReactionUnicodeEmoji    ReactionType = "unicode_emoji"
	ReactionRealmEmoji      ReactionType = "realm_emoji"
	ReactionZulipExtraEmoji ReactionType = "zulip_extra_emoji"
)

type OrganizationRoleLevel int

const (
	RoleOwner         OrganizationRoleLevel = 100
	RoleAdministrator OrganizationRoleLevel = 200
	RoleModerator     OrganizationRoleLevel = 300
	RoleMember        OrganizationRoleLevel = 400
	RoleGuest         OrganizationRoleLevel = 600
)
