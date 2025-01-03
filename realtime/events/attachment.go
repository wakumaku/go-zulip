package events

const AttachmentType EventType = "attachment"

type Attachment struct {
	ID              int       `json:"id"`
	Type            EventType `json:"type"`
	Op              string    `json:"op"`
	UploadSpaceUsed int       `json:"upload_space_used"`
	AttachmentData
}

type AttachmentData struct {
	Attachment struct {
		ID              int                 `json:"id"`
		Name            string              `json:"name"`
		PathID          string              `json:"path_id"`
		Size            int                 `json:"size"`
		CreateTime      int                 `json:"create_time"`
		Messages        []AttachmentMessage `json:"messages"`
		UploadSpaceUsed int                 `json:"upload_space_used"`
	} `json:"attachment"`
}

type AttachmentMessage struct {
	ID       int `json:"id"`
	DateSent int `json:"date_sent"`
}

func (e *Attachment) EventID() int {
	return e.ID
}

func (e *Attachment) EventType() EventType {
	return e.Type
}

func (e *Attachment) EventOp() string {
	return e.Op
}
