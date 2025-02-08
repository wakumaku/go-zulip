package events

const UserStatusType EventType = "user_status"

type UserStatus struct {
	ID   int       `json:"id"`
	Type EventType `json:"type"`
	UserStatusData
}

type UserStatusData struct {
	Away         bool   `json:"away"`
	EmojiCode    string `json:"emoji_code"`
	EmojiName    string `json:"emoji_name"`
	ReactionType string `json:"reaction_type"`
	StatusText   string `json:"status_text"`
	UserID       int    `json:"user_id"`
}

func (e *UserStatus) EventID() int {
	return e.ID
}

func (e *UserStatus) EventType() EventType {
	return e.Type
}

func (e *UserStatus) EventOp() string {
	return string(UserStatusType)
}
