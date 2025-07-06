package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/narrow"
	"github.com/wakumaku/go-zulip/realtime/events"
)

type RegisterEventQueueResponse struct {
	zulip.APIResponseBase
	registerEventQueueData
}

type registerEventQueueData struct {
	ZulipFeatureLevel int    `json:"zulip_feature_level"`
	ZulipMergeBase    string `json:"zulip_merge_base"`
	MaxMessageID      int    `json:"max_message_id"`
	LastEventID       int    `json:"last_event_id"`
	QueueID           string `json:"queue_id"`
	ZulipVersion      string `json:"zulip_version"`
}

func (r *RegisterEventQueueResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &r.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &r.registerEventQueueData); err != nil {
		return err
	}

	return nil
}

type ClientCapability string

const (
	NotificationSettingsNull   ClientCapability = "notification_settings_null"
	BulkMessageDeletion        ClientCapability = "bulk_message_deletion"
	UserAvatarURLFieldOptional ClientCapability = "user_avatar_url_field_optional"
	StreamTypingNotifications  ClientCapability = "stream_typing_notifications"
	UserSettingsObject         ClientCapability = "user_settings_object"
	LinkifierURLTemplate       ClientCapability = "linkifier_url_template"
	UserListIncomplete         ClientCapability = "user_list_incomplete"
	IncludeDeactivatedGroups   ClientCapability = "include_deactivated_groups"
	ArchivedChannels           ClientCapability = "archived_channels"
)

type registerEventQueueOptions struct {
	applyMarkdown            bool
	clientGravatar           *bool
	includeSubscribers       bool
	slimPresence             bool
	presenceHistoryLimitDays *int
	eventTypes               []events.EventType
	allPublicStreams         bool
	clientCapabilities       map[ClientCapability]bool
	fetchEventTypes          []events.EventType
	narrow                   narrow.Filter
}

type RegisterEventQueueOption func(*registerEventQueueOptions)

func ApplyMarkdown(apply bool) RegisterEventQueueOption {
	return func(ro *registerEventQueueOptions) {
		ro.applyMarkdown = apply
	}
}

func ClientGravatarEvent(clientGravatar bool) RegisterEventQueueOption {
	return func(ro *registerEventQueueOptions) {
		ro.clientGravatar = &clientGravatar
	}
}

func IncludeSubscribers(includeSubscribers bool) RegisterEventQueueOption {
	return func(ro *registerEventQueueOptions) {
		ro.includeSubscribers = includeSubscribers
	}
}

func SlimPresence(slimPresence bool) RegisterEventQueueOption {
	return func(ro *registerEventQueueOptions) {
		ro.slimPresence = slimPresence
	}
}

func PresenceHistoryLimitDays(presenceHistoryLimitDays int) RegisterEventQueueOption {
	return func(ro *registerEventQueueOptions) {
		ro.presenceHistoryLimitDays = &presenceHistoryLimitDays
	}
}

func EventTypes(eventTypes ...events.EventType) RegisterEventQueueOption {
	return func(ro *registerEventQueueOptions) {
		ro.eventTypes = eventTypes
	}
}

func AllPublicStreams(allPublicStreams bool) RegisterEventQueueOption {
	return func(ro *registerEventQueueOptions) {
		ro.allPublicStreams = allPublicStreams
	}
}

func ClientCapabilities(clientCapabilities map[ClientCapability]bool) RegisterEventQueueOption {
	return func(ro *registerEventQueueOptions) {
		ro.clientCapabilities = clientCapabilities
	}
}

func FetchEventTypes(fetchEventTypes []events.EventType) RegisterEventQueueOption {
	return func(ro *registerEventQueueOptions) {
		ro.fetchEventTypes = fetchEventTypes
	}
}

func NarrowEvents(narrow narrow.Filter) RegisterEventQueueOption {
	return func(ro *registerEventQueueOptions) {
		ro.narrow = narrow
	}
}

func (svc *Service) RegisterEvetQueue(ctx context.Context, options ...RegisterEventQueueOption) (*RegisterEventQueueResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/register"
	)

	// default options
	opts := registerEventQueueOptions{
		applyMarkdown:            false,
		clientGravatar:           nil,
		includeSubscribers:       false,
		slimPresence:             false,
		presenceHistoryLimitDays: nil,
		eventTypes:               []events.EventType{},
		allPublicStreams:         false,
		clientCapabilities:       map[ClientCapability]bool{},
		fetchEventTypes:          []events.EventType{},
		narrow:                   narrow.NewFilter(),
	}

	for _, opt := range options {
		opt(&opts)
	}

	msg := map[string]any{}

	if opts.applyMarkdown {
		msg["apply_markdown"] = opts.applyMarkdown
	}

	if opts.clientGravatar != nil {
		msg["client_gravatar"] = opts.clientGravatar
	}

	if opts.includeSubscribers {
		msg["include_subscribers"] = opts.includeSubscribers
	}

	if opts.slimPresence {
		msg["slim_presence"] = opts.slimPresence
	}

	if opts.presenceHistoryLimitDays != nil {
		msg["presence_history_limit_days"] = opts.presenceHistoryLimitDays
	}

	if len(opts.eventTypes) > 0 {
		eventTypes, err := json.Marshal(opts.eventTypes)
		if err != nil {
			return nil, fmt.Errorf("marshaling event types: %w", err)
		}

		msg["event_types"] = string(eventTypes)
	}

	if opts.allPublicStreams {
		msg["all_public_streams"] = opts.allPublicStreams
	}

	if len(opts.clientCapabilities) > 0 {
		clientCapabilities, err := json.Marshal(opts.clientCapabilities)
		if err != nil {
			return nil, fmt.Errorf("marshaling client capabilities: %w", err)
		}

		msg["client_capabilities"] = string(clientCapabilities)
	}

	if len(opts.fetchEventTypes) > 0 {
		fetchEventTypes, err := json.Marshal(opts.fetchEventTypes)
		if err != nil {
			return nil, fmt.Errorf("marshaling fetch event types: %w", err)
		}

		msg["fetch_event_types"] = string(fetchEventTypes)
	}

	if len(opts.narrow) > 0 {
		narrowJSON, err := opts.narrow.MarshalEvent()
		if err != nil {
			return nil, fmt.Errorf("marshaling narrow: %w", err)
		}

		msg["narrow"] = string(narrowJSON)
	}

	resp := RegisterEventQueueResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
