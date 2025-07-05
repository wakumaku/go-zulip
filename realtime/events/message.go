package events

import (
	"encoding/json"
	"errors"
)

const MessageType EventType = "message"

type Message struct {
	ID      int         `json:"id"`
	Type    EventType   `json:"type"`
	Message MessageData `json:"message"`
	Flags   []string    `json:"flags"`
}

type MessageData struct {
	ID               int              `json:"id"`
	Type             string           `json:"type"`
	AvatarURL        string           `json:"avatar_url"`
	Client           string           `json:"client"`
	Content          string           `json:"content"`
	ContentType      string           `json:"content_type"`
	DisplayRecipient DisplayRecipient `json:"display_recipient"`
	IsMeMessage      bool             `json:"is_me_message"`
	Reactions        []Reaction       `json:"reactions"`
	RecipientID      int              `json:"recipient_id"`
	SenderEmail      string           `json:"sender_email"`
	SenderFullName   string           `json:"sender_full_name"`
	SenderID         int              `json:"sender_id"`
	SenderRealmStr   string           `json:"sender_realm_str"`
	StreamID         int              `json:"stream_id"`
	Subject          string           `json:"subject"`
	Submessages      []Submessage     `json:"submessages"`
	Timestamp        int              `json:"timestamp"`
	TopicLinks       []TopicLinks     `json:"topic_links"`
}

type DisplayRecipient struct {
	IsChannel bool
	Channel   string
	Users     []DisplayRecipientObject
}

type DisplayRecipientObject struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

func (dr *DisplayRecipient) UnmarshalJSON(b []byte) error {
	var displayRecipient string
	if err := json.Unmarshal(b, &displayRecipient); err == nil {
		dr.IsChannel = true
		dr.Channel = displayRecipient

		return nil
	}

	var displayRecipientObject []DisplayRecipientObject
	if err := json.Unmarshal(b, &displayRecipientObject); err == nil {
		dr.IsChannel = false
		dr.Users = displayRecipientObject

		return nil
	}

	return errors.New("failed to unmarshal DisplayRecipient")
}

type Reaction struct {
	EmojiName    string `json:"emoji_name"`
	EmojiCode    string `json:"emoji_code"`
	ReactionType string `json:"reaction_type"`
}

type TopicLinks struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func (e *Message) EventID() int {
	return e.ID
}

func (e *Message) EventType() EventType {
	return e.Type
}

func (e *Message) EventOp() string {
	return "message"
}
