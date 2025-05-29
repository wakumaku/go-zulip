package events

const DeleteMessageType EventType = "delete_message"

type DeleteMessage struct {
	ID          int       `json:"id"`
	Type        EventType `json:"type"`
	MessageType string    `json:"message_type"`
	MessageID   *int      `json:"message_id"`
	MessageIDs  []int     `json:"message_ids"`
	StreamID    *int      `json:"stream_id"`
	Topic       *string   `json:"topic"`
}

func (e *DeleteMessage) EventID() int {
	return e.ID
}

func (e *DeleteMessage) EventType() EventType {
	return e.Type
}

func (e *DeleteMessage) EventOp() string {
	return "delete_message"
}
