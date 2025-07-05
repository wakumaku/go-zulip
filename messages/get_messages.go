package messages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/narrow"
)

type GetMessagesResponse struct {
	zulip.APIResponseBase
	getMessagesResponseData
}

type getMessagesResponseData struct {
	Anchor         int       `json:"anchor"`
	FoundNewest    bool      `json:"found_newest"`
	FoundOldest    bool      `json:"found_oldest"`
	FoundAnchor    bool      `json:"found_anchor"`
	HistoryLimited bool      `json:"history_limited"`
	Messages       []Message `json:"messages"`
}

func (g *GetMessagesResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getMessagesResponseData); err != nil {
		return err
	}

	return nil
}

type Message struct {
	AvatarURL         string           `json:"avatar_url"`
	Client            string           `json:"client"`
	Content           string           `json:"content"`
	ContentType       string           `json:"content_type"`
	DisplayRecipient  DisplayRecipient `json:"display_recipient"`
	EditHistory       []struct{}       `json:"edit_history"`
	ID                int              `json:"id"`
	IsMeMessage       bool             `json:"is_me_message"`
	LastEditTimestamp int              `json:"last_edit_timestamp"`
	Reactions         []struct {
		EmojiName    string           `json:"emoji_name"`
		EmojiCode    string           `json:"emoji_code"`
		ReactionType string           `json:"reaction_type"`
		UserID       int              `json:"user_id"`
		User         DisplayRecipient `json:"user"`
	} `json:"reactions"`
	RecipientID    int    `json:"recipient_id"`
	SenderEmail    string `json:"sender_email"`
	SenderFullName string `json:"sender_full_name"`
	SenderID       int    `json:"sender_id"`
	SenderRealmStr string `json:"sender_realm_str"`
	StreamID       int    `json:"stream_id"`
	Subject        string `json:"subject"`
	Submessages    []struct {
		MsgType   string `json:"msg_type"`
		Content   string `json:"content"`
		MessageID int    `json:"message_id"`
		SenderID  int    `json:"sender_id"`
		ID        int    `json:"id"`
	} `json:"submessages"`
	Timestamp  int `json:"timestamp"`
	TopicLinks []struct {
		Text string `json:"text"`
		URL  string `json:"url"`
	} `json:"topic_links"`
	Type         string   `json:"type"`
	Flags        []string `json:"flags"`
	MatchContent string   `json:"match_content"`
	MathSubject  string   `json:"match_subject"`
}

type DisplayRecipient struct {
	IsChannel bool
	Channel   string
	Users     []DisplayRecipientObject
}

type DisplayRecipientObject struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	FullName      string `json:"full_name"`
	IsMirrorDummy bool   `json:"is_mirror_dummy"`
}

func (d *DisplayRecipient) UnmarshalJSON(b []byte) error {
	var displayRecipient string
	if err := json.Unmarshal(b, &displayRecipient); err == nil {
		d.IsChannel = true
		d.Channel = displayRecipient

		return nil
	}

	var displayRecipientObject []DisplayRecipientObject
	if err := json.Unmarshal(b, &displayRecipientObject); err == nil {
		d.IsChannel = false
		d.Users = displayRecipientObject

		return nil
	}

	return errors.New("failed to unmarshal DisplayRecipient")
}

type getMessageOptions struct {
	anchor struct {
		fieldName string
		value     *string
	}
	includeAnchor struct {
		fieldName string
		value     *bool
	}
	numBefore struct {
		fieldName string
		value     *int
	}
	numAfter struct {
		fieldName string
		value     *int
	}
	narrow struct {
		fieldName string
		value     narrow.Filter
	}
	clientGravatar struct {
		fieldName string
		value     *bool
	}
	applyMarkdown struct {
		fieldName string
		value     *bool
	}
	messageIDs struct {
		fieldName string
		value     []int
	}

	// deprecated: useFirstUnreadAnchor Legacy way to specify "anchor": "first_unread" in Zulip 2.1.x and older.
	// useFirstUnreadAnchor struct {
	// 	fieldName string
	// 	value     *bool
	// }
}

type GetMessageOption func(*getMessageOptions)

func Anchor(anchor string) GetMessageOption {
	return func(o *getMessageOptions) {
		o.anchor.fieldName = "anchor"
		o.anchor.value = &anchor
	}
}

func IncludeAnchor(includeAnchor bool) GetMessageOption {
	return func(o *getMessageOptions) {
		o.includeAnchor.fieldName = "include_anchor"
		o.includeAnchor.value = &includeAnchor
	}
}

func NumBefore(numBefore int) GetMessageOption {
	return func(o *getMessageOptions) {
		o.numBefore.fieldName = "num_before"
		o.numBefore.value = &numBefore
	}
}

func NumAfter(numAfter int) GetMessageOption {
	return func(o *getMessageOptions) {
		o.numAfter.fieldName = "num_after"
		o.numAfter.value = &numAfter
	}
}

func NarrowMessage(narrow narrow.Filter) GetMessageOption {
	return func(o *getMessageOptions) {
		o.narrow.fieldName = "narrow"
		o.narrow.value = narrow
	}
}

func ClientGravatarMessage(clientGravatar bool) GetMessageOption {
	return func(o *getMessageOptions) {
		o.clientGravatar.fieldName = "client_gravatar"
		o.clientGravatar.value = &clientGravatar
	}
}

func ApplyMarkdownMessage(applyMarkdown bool) GetMessageOption {
	return func(o *getMessageOptions) {
		o.applyMarkdown.fieldName = "apply_markdown"
		o.applyMarkdown.value = &applyMarkdown
	}
}

func MessageIDs(messageIDs []int) GetMessageOption {
	return func(o *getMessageOptions) {
		o.messageIDs.fieldName = "message_ids"
		o.messageIDs.value = messageIDs
	}
}

func (svc *Service) GetMessages(ctx context.Context, options ...GetMessageOption) (*GetMessagesResponse, error) {
	const (
		method = http.MethodGet
		path   = "/api/v1/messages"
	)

	msg := map[string]any{}

	opts := getMessageOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.anchor.value != nil {
		msg[opts.anchor.fieldName] = *opts.anchor.value
	}

	if opts.includeAnchor.value != nil {
		msg[opts.includeAnchor.fieldName] = *opts.includeAnchor.value
	}

	if opts.numBefore.value != nil {
		msg[opts.numBefore.fieldName] = *opts.numBefore.value
	}

	if opts.numAfter.value != nil {
		msg[opts.numAfter.fieldName] = *opts.numAfter.value
	}

	if len(opts.narrow.value) > 0 {
		narrowJSON, err := opts.narrow.value.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("marshaling narrow: %w", err)
		}

		msg[opts.narrow.fieldName] = string(narrowJSON)
	}

	if opts.clientGravatar.value != nil {
		msg[opts.clientGravatar.fieldName] = *opts.clientGravatar.value
	}

	if opts.applyMarkdown.value != nil {
		msg[opts.applyMarkdown.fieldName] = *opts.applyMarkdown.value
	}

	if opts.messageIDs.value != nil {
		msg[opts.messageIDs.fieldName] = opts.messageIDs.value
	}

	resp := GetMessagesResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
