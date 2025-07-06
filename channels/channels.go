// Package channels provides functionality for managing Zulip channels (streams).
//
// Implemented features:
//   - Get subscribed channels
//   - Subscribe to channels (creates channels if they don't exist)
//   - Unsubscribe from channels
//   - Get all channels
//   - Get channel by ID
//   - Get channel ID
//   - Get channel subscribers
//   - Get subscription status
//
// See https://zulip.com/api/ for the complete API documentation.
package channels

import "github.com/wakumaku/go-zulip"

type Service struct {
	client zulip.RESTClient
}

func NewService(c zulip.RESTClient) *Service {
	return &Service{client: c}
}

// ChannelInfo are the fields usually returned when querying channel information
type ChannelInfo struct {
	CanAddSubscribersGroup     int    `json:"can_add_subscribers_group"`     // 10,
	CanRemoveSubscribersGroup  int    `json:"can_remove_subscribers_group"`  // 10,
	CreatorID                  int    `json:"creator_id"`                    // null,
	DateCreated                int    `json:"date_created"`                  // 1691057093,
	Description                string `json:"description"`                   // "A private channel",
	FirstMessageID             int    `json:"first_message_id"`              // 18,
	HistoryPublicToSubscribers bool   `json:"history_public_to_subscribers"` // false,
	InviteOnly                 bool   `json:"invite_only"`                   // true,
	IsAnnouncementOnly         bool   `json:"is_announcement_only"`          // false,
	IsArchived                 bool   `json:"is_archived"`                   // false,
	IsDefault                  bool   `json:"is_default"`                    // false,
	IsRecentlyActive           bool   `json:"is_recently_active"`            // true,
	IsWebPublic                bool   `json:"is_web_public"`                 // false,
	MessageRetentionDays       int    `json:"message_retention_days"`        // null,
	Name                       string `json:"name"`                          // "management",
	RenderedDescription        string `json:"rendered_description"`          // "<p>A private channel</p>",
	StreamID                   int    `json:"stream_id"`                     // 2,
	StreamPostPolicy           int    `json:"stream_post_policy"`            // 1,
	StreamWeeklyTraffic        int    `json:"stream_weekly_traffic"`         // null
}
