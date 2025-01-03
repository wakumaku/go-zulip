package events

const SubmessageType EventType = "submessage"

type Submessage struct {
	ID           int       `json:"id"`
	Type         EventType `json:"type"`
	Content      string    `json:"content"`
	MessageID    int       `json:"message_id"`
	MsgType      string    `json:"msg_type"`
	SenderID     int       `json:"sender_id"`
	SubmessageID int       `json:"submessage_id"`
}

func (e *Submessage) EventID() int {
	return e.ID
}

func (e *Submessage) EventType() EventType {
	return e.Type
}

func (e *Submessage) EventOp() string {
	return "submessage"
}
