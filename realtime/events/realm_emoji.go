package events

const RealmEmojiType EventType = "realm_emoji"

type RealmEmoji struct {
	ID   int       `json:"id"`
	Op   string    `json:"op"`
	Type EventType `json:"type"`
	RealmEmojiData
}

type RealmEmojiData struct {
	RealmEmoji map[string]struct {
		AuthorID    int    `json:"author_id"`
		Deactivated bool   `json:"deactivated"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		SourceURL   string `json:"source_url"`
	} `json:"realm_emoji"`
}

func (e *RealmEmoji) EventID() int {
	return e.ID
}

func (e *RealmEmoji) EventType() EventType {
	return e.Type
}

func (e *RealmEmoji) EventOp() string {
	return e.Op
}
