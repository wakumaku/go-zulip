package events

const HeartbeatType EventType = "heartbeat"

type Heartbeat struct {
	ID   int       `json:"id"`
	Type EventType `json:"type"`
}

func (e *Heartbeat) EventID() int {
	return e.ID
}

func (e *Heartbeat) EventType() EventType {
	return e.Type
}

func (e *Heartbeat) EventOp() string {
	return "heartbeat"
}
