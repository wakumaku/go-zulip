package events

const AlertWordsType EventType = "alert_words"

type AlertWords struct {
	ID   int       `json:"id"`
	Type EventType `json:"type"`
	AlertWordsData
}

type AlertWordsData struct {
	AlertWords []string `json:"alert_words"`
}

func (e *AlertWords) EventID() int {
	return e.ID
}

func (e *AlertWords) EventType() EventType {
	return e.Type
}

func (e *AlertWords) EventOp() string {
	return "alert_words"
}
