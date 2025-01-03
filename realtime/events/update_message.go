package events

const UpdateMessageType EventType = "update_message"

type UpdateMessage struct {
	ID                  int         `json:"id"`
	Type                EventType   `json:"type"`
	UserID              *int        `json:"user_id"`
	RenderingOnly       bool        `json:"rendering_only"`
	MessageID           int         `json:"message_id"`
	MessageIDs          []int       `json:"message_ids"`
	Flags               []string    `json:"flags"`
	EditTimestamp       int         `json:"edit_timestamp"`
	StreamName          *string     `json:"stream_name"`
	StreamID            *int        `json:"stream_id"`
	NewStreamID         *int        `json:"new_stream_id"`
	PropagateMode       *string     `json:"propagate_mode"`
	OrigSubject         *string     `json:"orig_subject"`
	Subject             *string     `json:"subject"`
	TopicLinks          []TopicLink `json:"topic_links"`
	OrigContent         *string     `json:"orig_content"`
	OrigRenderedContent *string     `json:"orig_rendered_content"`
	Content             *string     `json:"content"`
	RenderedContent     *string     `json:"rendered_content"`
	IsMeMessage         *bool       `json:"is_me_message"`
}

type TopicLink struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func (e *UpdateMessage) EventID() int {
	return e.ID
}

func (e *UpdateMessage) EventType() EventType {
	return e.Type
}

func (e *UpdateMessage) EventOp() string {
	return "update_message"
}
