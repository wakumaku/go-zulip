package events

const TypingType EventType = "typing"

type Typing struct {
	ID          int         `json:"id"`
	MessageType string      `json:"message_type"`
	Op          string      `json:"op"`
	Recipients  []Recipient `json:"recipients"`
	Sender      Sender      `json:"sender"`
	Type        EventType   `json:"type"`
}

type Recipient struct {
	Email  string `json:"email"`
	UserID int    `json:"user_id"`
}

type Sender struct {
	Email  string `json:"email"`
	UserID int    `json:"user_id"`
}

func (e *Typing) EventID() int {
	return e.ID
}

func (e *Typing) EventType() EventType {
	return e.Type
}

func (e *Typing) EventOp() string {
	return e.Op
}
