package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type UpdateSettingsResponse struct {
	zulip.APIResponseBase
}

func (g *UpdateSettingsResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	return nil
}

type webMarkReadOnScrollPolicy int

const (
	WebMarkReadOnScrollPolicyAlways                  webMarkReadOnScrollPolicy = 1
	WebMarkReadOnScrollPolicyOnlyInConversationViews webMarkReadOnScrollPolicy = 2
	WebMarkReadOnScrollPolicyNever                   webMarkReadOnScrollPolicy = 3
)

type webChannelDefaultView int

const (
	WebChannelDefaultViewTopTopicInTheChannel webChannelDefaultView = 1
	WebChannelDefaultViewChannelFeed          webChannelDefaultView = 2
)

type colorScheme int

const (
	ColorSchemeAutomatic  colorScheme = 1
	ColorSchemeDarkTheme  colorScheme = 2
	ColorSchemeLightTheme colorScheme = 3
)

type webHomeView string

const (
	WebHomeViewRecentTopics webHomeView = "recent_topics"
	WebHomeViewInbox        webHomeView = "inbox"
	WebHomeViewAllMessages  webHomeView = "all_messages"
)

type emojiset string

const (
	EmojiSetGoogle     emojiset = "google"
	EmojiSetGoogleBlob emojiset = "google-blob"
	EmojiSetTwitter    emojiset = "twitter"
	EmojiSetText       emojiset = "text"
)

type demoteInactiveStreams int

const (
	DemoteInactiveStreamsAutomatic demoteInactiveStreams = 1
	DemoteInactiveStreamsAlways    demoteInactiveStreams = 2
	DemoteInactiveStreamsNever     demoteInactiveStreams = 3
)

type userListStyle int

const (
	UserListStyleCompact             userListStyle = 1
	UserListStyleWithStatus          userListStyle = 2
	UserListStyleWithAvatarAndStatus userListStyle = 3
)

type webAnimateImagePreviews string

const (
	WebAnimateImagePreviewsAlways  webAnimateImagePreviews = "always"
	WebAnimateImagePreviewsOnHover webAnimateImagePreviews = "on_hover"
	WebAnimateImagePreviewsNever   webAnimateImagePreviews = "never"
)

type webStreamUnreadsCountDisplayPolicy int

const (
	WebStreamUnreadsCountDisplayPolicyAllChannels              webStreamUnreadsCountDisplayPolicy = 1
	WebStreamUnreadsCountDisplayPolicyUnmutedChannelsAndTopics webStreamUnreadsCountDisplayPolicy = 2
	WebStreamUnreadsCountDisplayPolicyNoChannels               webStreamUnreadsCountDisplayPolicy = 3
)

type desktopIconCountDisplay int

const (
	DesktopIconCountDisplayAllUnreadMessages            desktopIconCountDisplay = 1
	DesktopIconCountDisplayDMsMentionsAndFollowedTopics desktopIconCountDisplay = 2
	DesktopIconCountDisplayDMsAndMentions               desktopIconCountDisplay = 3
	DesktopIconCountDisplayNone                         desktopIconCountDisplay = 4
)

type realmNameInEmailNotificationsPolicy int

const (
	RealmNameInEmailNotificationsPolicyAutomatic realmNameInEmailNotificationsPolicy = 1
	RealmNameInEmailNotificationsPolicyAlways    realmNameInEmailNotificationsPolicy = 2
	RealmNameInEmailNotificationsPolicyNever     realmNameInEmailNotificationsPolicy = 3
)

type automaticallyFollowTopicsPolicy int

const (
	AutomaticallyFollowTopicsPolicyTopicsTheUserParticipatesIn  automaticallyFollowTopicsPolicy = 1
	AutomaticallyFollowTopicsPolicyTopicsTheUserSendsAMessageTo automaticallyFollowTopicsPolicy = 2
	AutomaticallyFollowTopicsPolicyTopicsTheUserStarts          automaticallyFollowTopicsPolicy = 3
	AutomaticallyFollowTopicsPolicyNever                        automaticallyFollowTopicsPolicy = 4
)

type emailAddressVisibility int

const (
	EmailAddressVisibilityEveryone           emailAddressVisibility = 1
	EmailAddressVisibilityMembersOnly        emailAddressVisibility = 2
	EmailAddressVisibilityAdministratorsOnly emailAddressVisibility = 3
	EmailAddressVisibilityNobody             emailAddressVisibility = 4
	EmailAddressVisibilityModeratorsOnly     emailAddressVisibility = 5
)

type updateSettingsOptions struct {
	// email, string
	email struct {
		fieldName string
		value     *string
	}
	// full_name, string
	fullName struct {
		fieldName string
		value     *string
	}
	// old_password, string
	oldPassword struct {
		fieldName string
		value     *string
	}
	// new_password, string (old_password, string)
	newPassword struct {
		fieldName string
		value     *string
	}
	// twenty_four_hour_time, boolean
	twentyFourHourTime struct {
		fieldName string
		value     *bool
	}
	// web_mark_read_on_scroll_policy, integer [1: Always, 2: Only in conversation views, 3: Never]
	webMarkReadOnScrollPolicy struct {
		fieldName string
		value     *webMarkReadOnScrollPolicy
	}
	// web_channel_default_view, integer [1: Top topic in the channel, 2: Channel feed]
	webChannelDefaultView struct {
		fieldName string
		value     *webChannelDefaultView
	}
	// starred_message_counts, boolean
	starredMessageCounts struct {
		fieldName string
		value     *bool
	}
	// receives_typing_notifications, boolean
	receivesTypingNotifications struct {
		fieldName string
		value     *bool
	}
	// web_suggest_update_timezone, boolean
	webSuggestUpdateTimezone struct {
		fieldName string
		value     *bool
	}
	// fluid_layout_width, boolean
	fluidLayoutWidth struct {
		fieldName string
		value     *bool
	}
	// high_contrast_mode, boolean
	highContrastMode struct {
		fieldName string
		value     *bool
	}
	// web_font_size_px, integer
	webFontSizePx struct {
		fieldName string
		value     *int
	}
	// web_line_height_percent, integer
	webLineHeightPercent struct {
		fieldName string
		value     *int
	}
	// color_scheme, integer [1: Automatic, 2: Dark theme, 3: Light theme]
	colorScheme struct {
		fieldName string
		value     *colorScheme
	}

	// enable_drafts_synchronization, boolean
	enableDraftsSynchronization struct {
		fieldName string
		value     *bool
	}

	// translate_emoticons, boolean
	translateEmoticons struct {
		fieldName string
		value     *bool
	}

	// display_emoji_reaction_users, boolean
	displayEmojiReactionUsers struct {
		fieldName string
		value     *bool
	}

	// default_language, string
	defaultLanguage struct {
		fieldName string
		value     *string
	}

	// web_home_view, string ["recent_topics", "inbox", "all_messages"]
	webHomeView struct {
		fieldName string
		value     *webHomeView
	}

	// web_escape_navigates_to_home_view, boolean
	webEscapeNavigatesToHomeView struct {
		fieldName string
		value     *bool
	}

	// left_side_userlist, boolean
	leftSideUserlist struct {
		fieldName string
		value     *bool
	}

	// emojiset, string ["google", "google-blob", "twitter", "text"]
	emojiset struct {
		fieldName string
		value     *emojiset
	}

	// demote_inactive_streams, integer [1: Automatic, 2: Always, 3: Never]
	demoteInactiveStreams struct {
		fieldName string
		value     *demoteInactiveStreams
	}

	// user_list_style, integer [1: Compact, 2: With status, 3: With avatar and status]
	userListStyle struct {
		fieldName string
		value     *userListStyle
	}

	// web_animate_image_previews, string ["always", "on_hover", "never"]
	webAnimateImagePreviews struct {
		fieldName string
		value     *webAnimateImagePreviews
	}

	// web_stream_unreads_count_display_policy, integer [1: All channels, 2: Unmuted channels and topics, 3: No channels]
	webStreamUnreadsCountDisplayPolicy struct {
		fieldName string
		value     *webStreamUnreadsCountDisplayPolicy
	}

	// hide_ai_features, boolean
	hideAiFeatures struct {
		fieldName string
		value     *bool
	}

	// timezone, string
	timezone struct {
		fieldName string
		value     *string
	}

	// enable_stream_desktop_notifications, boolean
	enableStreamDesktopNotifications struct {
		fieldName string
		value     *bool
	}

	// enable_stream_email_notifications, boolean
	enableStreamEmailNotifications struct {
		fieldName string
		value     *bool
	}
	// enable_stream_push_notifications, boolean
	enableStreamPushNotifications struct {
		fieldName string
		value     *bool
	}
	// enable_stream_audible_notifications, boolean
	enableStreamAudibleNotifications struct {
		fieldName string
		value     *bool
	}
	// notification_sound, string
	notificationSound struct {
		fieldName string
		value     *string
	}
	// enable_desktop_notifications, boolean
	enableDesktopNotifications struct {
		fieldName string
		value     *bool
	}
	// enable_sounds, boolean
	enableSounds struct {
		fieldName string
		value     *bool
	}
	// email_notifications_batching_period_seconds, integer
	emailNotificationsBatchingPeriodSeconds struct {
		fieldName string
		value     *int
	}
	// enable_offline_email_notifications, boolean
	enableOfflineEmailNotifications struct {
		fieldName string
		value     *bool
	}
	// enable_offline_push_notifications, boolean
	enableOfflinePushNotifications struct {
		fieldName string
		value     *bool
	}
	// enable_online_push_notifications, boolean
	enableOnlinePushNotifications struct {
		fieldName string
		value     *bool
	}
	// enable_followed_topic_desktop_notifications, boolean
	enableFollowedTopicDesktopNotifications struct {
		fieldName string
		value     *bool
	}
	// enable_followed_topic_email_notifications, boolean
	enableFollowedTopicEmailNotifications struct {
		fieldName string
		value     *bool
	}
	// enable_followed_topic_push_notifications, boolean
	enableFollowedTopicPushNotifications struct {
		fieldName string
		value     *bool
	}
	// enable_followed_topic_audible_notifications, boolean
	enableFollowedTopicAudibleNotifications struct {
		fieldName string
		value     *bool
	}
	// enable_digest_emails, boolean
	enableDigestEmails struct {
		fieldName string
		value     *bool
	}
	// enable_marketing_emails, boolean
	enableMarketingEmails struct {
		fieldName string
		value     *bool
	}
	// enable_login_emails, boolean
	enableLoginEmails struct {
		fieldName string
		value     *bool
	}
	// message_content_in_email_notifications, boolean
	messageContentInEmailNotifications struct {
		fieldName string
		value     *bool
	}
	// pm_content_in_desktop_notifications, boolean
	pmContentInDesktopNotifications struct {
		fieldName string
		value     *bool
	}
	// wildcard_mentions_notify, boolean
	wildcardMentionsNotify struct {
		fieldName string
		value     *bool
	}
	// enable_followed_topic_wildcard_mentions_notify, boolean
	enableFollowedTopicWildcardMentionsNotify struct {
		fieldName string
		value     *bool
	}
	// desktop_icon_count_display, integer [1: All unread messages, 2: DMs, mentions, and followed topics, 3: DMs and mentions, 4: None]
	desktopIconCountDisplay struct {
		fieldName string
		value     *desktopIconCountDisplay
	}
	// realm_name_in_email_notifications_policy, integer [1: Automatic, 2: Always, 3: Never]
	realmNameInEmailNotificationsPolicy struct {
		fieldName string
		value     *realmNameInEmailNotificationsPolicy
	}
	// automatically_follow_topics_policy, integer [1: Topics the user participates in, 2: Topics the user sends a message to, 3: Topics the user starts, 4: Never]
	automaticallyFollowTopicsPolicy struct {
		fieldName string
		value     *automaticallyFollowTopicsPolicy
	}
	// automatically_unmute_topics_in_muted_streams_policy, integer [1: Topics the user participates in, 2: Topics the user sends a message to, 3: Topics the user starts, 4: Never]
	automaticallyUnmuteTopicsInMutedStreamsPolicy struct {
		fieldName string
		value     *automaticallyFollowTopicsPolicy
	}
	// automatically_follow_topics_where_mentioned, boolean
	automaticallyFollowTopicsWhereMentioned struct {
		fieldName string
		value     *bool
	}
	// presence_enabled, boolean
	presenceEnabled struct {
		fieldName string
		value     *bool
	}
	// enter_sends, boolean
	enterSends struct {
		fieldName string
		value     *bool
	}
	// send_private_typing_notifications, boolean
	sendPrivateTypingNotifications struct {
		fieldName string
		value     *bool
	}
	// send_stream_typing_notifications, boolean
	sendStreamTypingNotifications struct {
		fieldName string
		value     *bool
	}
	// send_read_receipts, boolean
	sendReadReceipts struct {
		fieldName string
		value     *bool
	}
	// allow_private_data_export, boolean
	allowPrivateDataExport struct {
		fieldName string
		value     *bool
	}
	// email_address_visibility, integer [1: Everyone, 2: Members only, 3: Administrators only, 4: Nobody, 5: Moderators only]
	emailAddressVisibility struct {
		fieldName string
		value     *emailAddressVisibility
	}
	// web_navigate_to_sent_message, boolean
	webNavigateToSentMessage struct {
		fieldName string
		value     *bool
	}
}

type UpdateSettingsOption func(*updateSettingsOptions)

func Email(value string) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.email.fieldName = "email"
		args.email.value = &value
	}
}

func SetFullName(value string) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.fullName.fieldName = "full_name"
		args.fullName.value = &value
	}
}

func SetPassword(newPassword, oldPassword string) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.newPassword.fieldName = "new_password"
		args.newPassword.value = &newPassword
		args.oldPassword.fieldName = "old_password"
		args.oldPassword.value = &oldPassword
	}
}

func TwentyFourHourTime(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.twentyFourHourTime.fieldName = "twenty_four_hour_time"
		args.twentyFourHourTime.value = &value
	}
}

func WebMarkReadOnScrollPolicy(value webMarkReadOnScrollPolicy) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.webMarkReadOnScrollPolicy.fieldName = "web_mark_read_on_scroll_policy"
		args.webMarkReadOnScrollPolicy.value = &value
	}
}

func WebChannelDefaultView(value webChannelDefaultView) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.webChannelDefaultView.fieldName = "web_channel_default_view"
		args.webChannelDefaultView.value = &value
	}
}

func StarredMessageCounts(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.starredMessageCounts.fieldName = "starred_message_counts"
		args.starredMessageCounts.value = &value
	}
}

func ReceivesTypingNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.receivesTypingNotifications.fieldName = "receives_typing_notifications"
		args.receivesTypingNotifications.value = &value
	}
}

func WebSuggestUpdateTimezone(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.webSuggestUpdateTimezone.fieldName = "web_suggest_update_timezone"
		args.webSuggestUpdateTimezone.value = &value
	}
}

func FluidLayoutWidth(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.fluidLayoutWidth.fieldName = "fluid_layout_width"
		args.fluidLayoutWidth.value = &value
	}
}

func HighContrastMode(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.highContrastMode.fieldName = "high_contrast_mode"
		args.highContrastMode.value = &value
	}
}

func WebFontSizePx(value int) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.webFontSizePx.fieldName = "web_font_size_px"
		args.webFontSizePx.value = &value
	}
}

func WebLineHeightPercent(value int) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.webLineHeightPercent.fieldName = "web_line_height_percent"
		args.webLineHeightPercent.value = &value
	}
}

func ColorScheme(value colorScheme) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.colorScheme.fieldName = "color_scheme"
		args.colorScheme.value = &value
	}
}

func EnableDraftsSynchronization(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableDraftsSynchronization.fieldName = "enable_drafts_synchronization"
		args.enableDraftsSynchronization.value = &value
	}
}

func TranslateEmoticons(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.translateEmoticons.fieldName = "translate_emoticons"
		args.translateEmoticons.value = &value
	}
}

func DisplayEmojiReactionUsers(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.displayEmojiReactionUsers.fieldName = "display_emoji_reaction_users"
		args.displayEmojiReactionUsers.value = &value
	}
}

func DefaultLanguage(value string) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.defaultLanguage.fieldName = "default_language"
		args.defaultLanguage.value = &value
	}
}

func WebHomeView(value webHomeView) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.webHomeView.fieldName = "web_home_view"
		args.webHomeView.value = &value
	}
}

func WebEscapeNavigatesToHomeView(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.webEscapeNavigatesToHomeView.fieldName = "web_escape_navigates_to_home_view"
		args.webEscapeNavigatesToHomeView.value = &value
	}
}

func LeftSideUserlist(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.leftSideUserlist.fieldName = "left_side_userlist"
		args.leftSideUserlist.value = &value
	}
}

func Emojiset(value emojiset) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.emojiset.fieldName = "emojiset"
		args.emojiset.value = &value
	}
}

func DemoteInactiveStreams(value demoteInactiveStreams) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.demoteInactiveStreams.fieldName = "demote_inactive_streams"
		args.demoteInactiveStreams.value = &value
	}
}

func UserListStyle(value userListStyle) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.userListStyle.fieldName = "user_list_style"
		args.userListStyle.value = &value
	}
}

func WebAnimateImagePreviews(value webAnimateImagePreviews) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.webAnimateImagePreviews.fieldName = "web_animate_image_previews"
		args.webAnimateImagePreviews.value = &value
	}
}

func WebStreamUnreadsCountDisplayPolicy(value webStreamUnreadsCountDisplayPolicy) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.webStreamUnreadsCountDisplayPolicy.fieldName = "web_stream_unreads_count_display_policy"
		args.webStreamUnreadsCountDisplayPolicy.value = &value
	}
}

func HideAiFeatures(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.hideAiFeatures.fieldName = "hide_ai_features"
		args.hideAiFeatures.value = &value
	}
}

func Timezone(value string) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.timezone.fieldName = "timezone"
		args.timezone.value = &value
	}
}

func EnableStreamDesktopNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableStreamDesktopNotifications.fieldName = "enable_stream_desktop_notifications"
		args.enableStreamDesktopNotifications.value = &value
	}
}

func EnableStreamEmailNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableStreamEmailNotifications.fieldName = "enable_stream_email_notifications"
		args.enableStreamEmailNotifications.value = &value
	}
}

func EnableStreamPushNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableStreamPushNotifications.fieldName = "enable_stream_push_notifications"
		args.enableStreamPushNotifications.value = &value
	}
}

func EnableStreamAudibleNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableStreamAudibleNotifications.fieldName = "enable_stream_audible_notifications"
		args.enableStreamAudibleNotifications.value = &value
	}
}

func NotificationSound(value string) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.notificationSound.fieldName = "notification_sound"
		args.notificationSound.value = &value
	}
}

func EnableDesktopNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableDesktopNotifications.fieldName = "enable_desktop_notifications"
		args.enableDesktopNotifications.value = &value
	}
}

func EnableSounds(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableSounds.fieldName = "enable_sounds"
		args.enableSounds.value = &value
	}
}

func EmailNotificationsBatchingPeriodSeconds(value int) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.emailNotificationsBatchingPeriodSeconds.fieldName = "email_notifications_batching_period_seconds"
		args.emailNotificationsBatchingPeriodSeconds.value = &value
	}
}

func EnableOfflineEmailNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableOfflineEmailNotifications.fieldName = "enable_offline_email_notifications"
		args.enableOfflineEmailNotifications.value = &value
	}
}

func EnableOfflinePushNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableOfflinePushNotifications.fieldName = "enable_offline_push_notifications"
		args.enableOfflinePushNotifications.value = &value
	}
}

func EnableOnlinePushNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableOnlinePushNotifications.fieldName = "enable_online_push_notifications"
		args.enableOnlinePushNotifications.value = &value
	}
}

func EnableFollowedTopicDesktopNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableFollowedTopicDesktopNotifications.fieldName = "enable_followed_topic_desktop_notifications"
		args.enableFollowedTopicDesktopNotifications.value = &value
	}
}

func EnableFollowedTopicEmailNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableFollowedTopicEmailNotifications.fieldName = "enable_followed_topic_email_notifications"
		args.enableFollowedTopicEmailNotifications.value = &value
	}
}

func EnableFollowedTopicPushNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableFollowedTopicPushNotifications.fieldName = "enable_followed_topic_push_notifications"
		args.enableFollowedTopicPushNotifications.value = &value
	}
}

func EnableFollowedTopicAudibleNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableFollowedTopicAudibleNotifications.fieldName = "enable_followed_topic_audible_notifications"
		args.enableFollowedTopicAudibleNotifications.value = &value
	}
}

func EnableDigestEmails(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableDigestEmails.fieldName = "enable_digest_emails"
		args.enableDigestEmails.value = &value
	}
}

func EnableMarketingEmails(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableMarketingEmails.fieldName = "enable_marketing_emails"
		args.enableMarketingEmails.value = &value
	}
}

func EnableLoginEmails(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableLoginEmails.fieldName = "enable_login_emails"
		args.enableLoginEmails.value = &value
	}
}

func MessageContentInEmailNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.messageContentInEmailNotifications.fieldName = "message_content_in_email_notifications"
		args.messageContentInEmailNotifications.value = &value
	}
}

func PMContentInDesktopNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.pmContentInDesktopNotifications.fieldName = "pm_content_in_desktop_notifications"
		args.pmContentInDesktopNotifications.value = &value
	}
}

func WildcardMentionsNotify(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.wildcardMentionsNotify.fieldName = "wildcard_mentions_notify"
		args.wildcardMentionsNotify.value = &value
	}
}

func EnableFollowedTopicWildcardMentionsNotify(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enableFollowedTopicWildcardMentionsNotify.fieldName = "enable_followed_topic_wildcard_mentions_notify"
		args.enableFollowedTopicWildcardMentionsNotify.value = &value
	}
}

func DesktopIconCountDisplay(value desktopIconCountDisplay) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.desktopIconCountDisplay.fieldName = "desktop_icon_count_display"
		args.desktopIconCountDisplay.value = &value
	}
}

func RealmNameInEmailNotificationsPolicy(value realmNameInEmailNotificationsPolicy) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.realmNameInEmailNotificationsPolicy.fieldName = "realm_name_in_email_notifications_policy"
		args.realmNameInEmailNotificationsPolicy.value = &value
	}
}

func AutomaticallyFollowTopicsPolicy(value automaticallyFollowTopicsPolicy) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.automaticallyFollowTopicsPolicy.fieldName = "automatically_follow_topics_policy"
		args.automaticallyFollowTopicsPolicy.value = &value
	}
}

func AutomaticallyUnmuteTopicsInMutedStreamsPolicy(value automaticallyFollowTopicsPolicy) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.automaticallyUnmuteTopicsInMutedStreamsPolicy.fieldName = "automatically_unmute_topics_in_muted_streams_policy"
		args.automaticallyUnmuteTopicsInMutedStreamsPolicy.value = &value
	}
}

func AutomaticallyFollowTopicsWhereMentioned(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.automaticallyFollowTopicsWhereMentioned.fieldName = "automatically_follow_topics_where_mentioned"
		args.automaticallyFollowTopicsWhereMentioned.value = &value
	}
}

func PresenceEnabled(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.presenceEnabled.fieldName = "presence_enabled"
		args.presenceEnabled.value = &value
	}
}

func EnterSends(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.enterSends.fieldName = "enter_sends"
		args.enterSends.value = &value
	}
}

func SendPrivateTypingNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.sendPrivateTypingNotifications.fieldName = "send_private_typing_notifications"
		args.sendPrivateTypingNotifications.value = &value
	}
}

func SendStreamTypingNotifications(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.sendStreamTypingNotifications.fieldName = "send_stream_typing_notifications"
		args.sendStreamTypingNotifications.value = &value
	}
}

func SendReadReceipts(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.sendReadReceipts.fieldName = "send_read_receipts"
		args.sendReadReceipts.value = &value
	}
}

func AllowPrivateDataExport(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.allowPrivateDataExport.fieldName = "allow_private_data_export"
		args.allowPrivateDataExport.value = &value
	}
}

func EmailAddressVisibility(value emailAddressVisibility) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.emailAddressVisibility.fieldName = "email_address_visibility"
		args.emailAddressVisibility.value = &value
	}
}

func WebNavigateToSentMessage(value bool) UpdateSettingsOption {
	return func(args *updateSettingsOptions) {
		args.webNavigateToSentMessage.fieldName = "web_navigate_to_sent_message"
		args.webNavigateToSentMessage.value = &value
	}
}

func (svc *Service) UpdateSettings(ctx context.Context, opts ...UpdateSettingsOption) (*UpdateSettingsResponse, error) {
	const (
		path   = "/api/v1/settings"
		method = http.MethodPatch
	)

	msg := map[string]any{}

	options := &updateSettingsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.email.value != nil {
		msg[options.email.fieldName] = *options.email.value
	}
	if options.fullName.value != nil {
		msg[options.fullName.fieldName] = *options.fullName.value
	}
	if options.oldPassword.value != nil {
		msg[options.oldPassword.fieldName] = *options.oldPassword.value
	}
	if options.newPassword.value != nil {
		msg[options.newPassword.fieldName] = *options.newPassword.value
	}
	if options.twentyFourHourTime.value != nil {
		msg[options.twentyFourHourTime.fieldName] = *options.twentyFourHourTime.value
	}
	if options.webMarkReadOnScrollPolicy.value != nil {
		msg[options.webMarkReadOnScrollPolicy.fieldName] = *options.webMarkReadOnScrollPolicy.value
	}
	if options.webChannelDefaultView.value != nil {
		msg[options.webChannelDefaultView.fieldName] = *options.webChannelDefaultView.value
	}
	if options.starredMessageCounts.value != nil {
		msg[options.starredMessageCounts.fieldName] = *options.starredMessageCounts.value
	}
	if options.receivesTypingNotifications.value != nil {
		msg[options.receivesTypingNotifications.fieldName] = *options.receivesTypingNotifications.value
	}
	if options.webSuggestUpdateTimezone.value != nil {
		msg[options.webSuggestUpdateTimezone.fieldName] = *options.webSuggestUpdateTimezone.value
	}
	if options.fluidLayoutWidth.value != nil {
		msg[options.fluidLayoutWidth.fieldName] = *options.fluidLayoutWidth.value
	}
	if options.highContrastMode.value != nil {
		msg[options.highContrastMode.fieldName] = *options.highContrastMode.value
	}
	if options.webFontSizePx.value != nil {
		msg[options.webFontSizePx.fieldName] = *options.webFontSizePx.value
	}
	if options.webLineHeightPercent.value != nil {
		msg[options.webLineHeightPercent.fieldName] = *options.webLineHeightPercent.value
	}
	if options.colorScheme.value != nil {
		msg[options.colorScheme.fieldName] = *options.colorScheme.value
	}
	if options.enableDraftsSynchronization.value != nil {
		msg[options.enableDraftsSynchronization.fieldName] = *options.enableDraftsSynchronization.value
	}
	if options.translateEmoticons.value != nil {
		msg[options.translateEmoticons.fieldName] = *options.translateEmoticons.value
	}
	if options.displayEmojiReactionUsers.value != nil {
		msg[options.displayEmojiReactionUsers.fieldName] = *options.displayEmojiReactionUsers.value
	}
	if options.defaultLanguage.value != nil {
		msg[options.defaultLanguage.fieldName] = *options.defaultLanguage.value
	}
	if options.webHomeView.value != nil {
		msg[options.webHomeView.fieldName] = *options.webHomeView.value
	}
	if options.webEscapeNavigatesToHomeView.value != nil {
		msg[options.webEscapeNavigatesToHomeView.fieldName] = *options.webEscapeNavigatesToHomeView.value
	}
	if options.leftSideUserlist.value != nil {
		msg[options.leftSideUserlist.fieldName] = *options.leftSideUserlist.value
	}
	if options.emojiset.value != nil {
		msg[options.emojiset.fieldName] = *options.emojiset.value
	}
	if options.demoteInactiveStreams.value != nil {
		msg[options.demoteInactiveStreams.fieldName] = *options.demoteInactiveStreams.value
	}
	if options.userListStyle.value != nil {
		msg[options.userListStyle.fieldName] = *options.userListStyle.value
	}
	if options.webAnimateImagePreviews.value != nil {
		msg[options.webAnimateImagePreviews.fieldName] = *options.webAnimateImagePreviews.value
	}
	if options.webStreamUnreadsCountDisplayPolicy.value != nil {
		msg[options.webStreamUnreadsCountDisplayPolicy.fieldName] = *options.webStreamUnreadsCountDisplayPolicy.value
	}
	if options.hideAiFeatures.value != nil {
		msg[options.hideAiFeatures.fieldName] = *options.hideAiFeatures.value
	}
	if options.timezone.value != nil {
		msg[options.timezone.fieldName] = *options.timezone.value
	}
	if options.enableStreamDesktopNotifications.value != nil {
		msg[options.enableStreamDesktopNotifications.fieldName] = *options.enableStreamDesktopNotifications.value
	}
	if options.enableStreamEmailNotifications.value != nil {
		msg[options.enableStreamEmailNotifications.fieldName] = *options.enableStreamEmailNotifications.value
	}
	if options.enableStreamPushNotifications.value != nil {
		msg[options.enableStreamPushNotifications.fieldName] = *options.enableStreamPushNotifications.value
	}
	if options.enableStreamAudibleNotifications.value != nil {
		msg[options.enableStreamAudibleNotifications.fieldName] = *options.enableStreamAudibleNotifications.value
	}
	if options.notificationSound.value != nil {
		msg[options.notificationSound.fieldName] = *options.notificationSound.value
	}
	if options.enableDesktopNotifications.value != nil {
		msg[options.enableDesktopNotifications.fieldName] = *options.enableDesktopNotifications.value
	}
	if options.enableSounds.value != nil {
		msg[options.enableSounds.fieldName] = *options.enableSounds.value
	}
	if options.emailNotificationsBatchingPeriodSeconds.value != nil {
		msg[options.emailNotificationsBatchingPeriodSeconds.fieldName] = *options.emailNotificationsBatchingPeriodSeconds.value
	}
	if options.enableOfflineEmailNotifications.value != nil {
		msg[options.enableOfflineEmailNotifications.fieldName] = *options.enableOfflineEmailNotifications.value
	}
	if options.enableOfflinePushNotifications.value != nil {
		msg[options.enableOfflinePushNotifications.fieldName] = *options.enableOfflinePushNotifications.value
	}
	if options.enableOnlinePushNotifications.value != nil {
		msg[options.enableOnlinePushNotifications.fieldName] = *options.enableOnlinePushNotifications.value
	}
	if options.enableFollowedTopicDesktopNotifications.value != nil {
		msg[options.enableFollowedTopicDesktopNotifications.fieldName] = *options.enableFollowedTopicDesktopNotifications.value
	}
	if options.enableFollowedTopicEmailNotifications.value != nil {
		msg[options.enableFollowedTopicEmailNotifications.fieldName] = *options.enableFollowedTopicEmailNotifications.value
	}
	if options.enableFollowedTopicPushNotifications.value != nil {
		msg[options.enableFollowedTopicPushNotifications.fieldName] = *options.enableFollowedTopicPushNotifications.value
	}
	if options.enableFollowedTopicAudibleNotifications.value != nil {
		msg[options.enableFollowedTopicAudibleNotifications.fieldName] = *options.enableFollowedTopicAudibleNotifications.value
	}
	if options.enableDigestEmails.value != nil {
		msg[options.enableDigestEmails.fieldName] = *options.enableDigestEmails.value
	}
	if options.enableMarketingEmails.value != nil {
		msg[options.enableMarketingEmails.fieldName] = *options.enableMarketingEmails.value
	}
	if options.enableLoginEmails.value != nil {
		msg[options.enableLoginEmails.fieldName] = *options.enableLoginEmails.value
	}
	if options.messageContentInEmailNotifications.value != nil {
		msg[options.messageContentInEmailNotifications.fieldName] = *options.messageContentInEmailNotifications.value
	}
	if options.pmContentInDesktopNotifications.value != nil {
		msg[options.pmContentInDesktopNotifications.fieldName] = *options.pmContentInDesktopNotifications.value
	}
	if options.wildcardMentionsNotify.value != nil {
		msg[options.wildcardMentionsNotify.fieldName] = *options.wildcardMentionsNotify.value
	}
	if options.enableFollowedTopicWildcardMentionsNotify.value != nil {
		msg[options.enableFollowedTopicWildcardMentionsNotify.fieldName] = *options.enableFollowedTopicWildcardMentionsNotify.value
	}
	if options.desktopIconCountDisplay.value != nil {
		msg[options.desktopIconCountDisplay.fieldName] = *options.desktopIconCountDisplay.value
	}
	if options.realmNameInEmailNotificationsPolicy.value != nil {
		msg[options.realmNameInEmailNotificationsPolicy.fieldName] = *options.realmNameInEmailNotificationsPolicy.value
	}
	if options.automaticallyFollowTopicsPolicy.value != nil {
		msg[options.automaticallyFollowTopicsPolicy.fieldName] = *options.automaticallyFollowTopicsPolicy.value
	}
	if options.automaticallyUnmuteTopicsInMutedStreamsPolicy.value != nil {
		msg[options.automaticallyUnmuteTopicsInMutedStreamsPolicy.fieldName] = *options.automaticallyUnmuteTopicsInMutedStreamsPolicy.value
	}
	if options.automaticallyFollowTopicsWhereMentioned.value != nil {
		msg[options.automaticallyFollowTopicsWhereMentioned.fieldName] = *options.automaticallyFollowTopicsWhereMentioned.value
	}
	if options.presenceEnabled.value != nil {
		msg[options.presenceEnabled.fieldName] = *options.presenceEnabled.value
	}
	if options.enterSends.value != nil {
		msg[options.enterSends.fieldName] = *options.enterSends.value
	}
	if options.sendPrivateTypingNotifications.value != nil {
		msg[options.sendPrivateTypingNotifications.fieldName] = *options.sendPrivateTypingNotifications.value
	}
	if options.sendStreamTypingNotifications.value != nil {
		msg[options.sendStreamTypingNotifications.fieldName] = *options.sendStreamTypingNotifications.value
	}
	if options.sendReadReceipts.value != nil {
		msg[options.sendReadReceipts.fieldName] = *options.sendReadReceipts.value
	}
	if options.allowPrivateDataExport.value != nil {
		msg[options.allowPrivateDataExport.fieldName] = *options.allowPrivateDataExport.value
	}
	if options.emailAddressVisibility.value != nil {
		msg[options.emailAddressVisibility.fieldName] = *options.emailAddressVisibility.value
	}
	if options.webNavigateToSentMessage.value != nil {
		msg[options.webNavigateToSentMessage.fieldName] = *options.webNavigateToSentMessage.value
	}

	resp := UpdateSettingsResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
