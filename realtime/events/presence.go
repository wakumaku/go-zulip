package events

const PresenceType EventType = "presence"

type Presence struct {
	Email           string       `json:"email"`
	ID              int          `json:"id"`
	Presence        PresenceData `json:"presence"`
	ServerTimestamp float64      `json:"server_timestamp"`
	Type            EventType    `json:"type"`
	UserID          int          `json:"user_id"`
}

type PresenceData struct {
	Website PresenceDetail `json:"website"`
}

type PresenceDetail struct {
	Client    string `json:"client"`
	Pushable  bool   `json:"pushable"`
	Status    string `json:"status"`
	Timestamp int    `json:"timestamp"`
}

func (e *Presence) EventID() int {
	return e.ID
}

func (e *Presence) EventType() EventType {
	return e.Type
}

func (e *Presence) EventOp() string {
	return "presence"
}
