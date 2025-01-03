package events

const UnknownType EventType = "unknown"

type Unknown map[string]any

func (e *Unknown) EventID() int {
	return -1
}

func (e *Unknown) EventType() EventType {
	return "unknown"
}

func (e *Unknown) EventOp() string {
	return "unknown"
}
