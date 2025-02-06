package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestAlertWords(t *testing.T) {
	eventExample := `{
    "alert_words": [
        "alert_word"
    ],
    "id": 0,
    "type": "alert_words"
}`

	v := events.AlertWords{}
	err := json.Unmarshal([]byte(eventExample), &v)
	assert.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.AlertWordsType, v.EventType())
	assert.Equal(t, "alert_words", v.EventOp())
	assert.Equal(t, []string{"alert_word"}, v.AlertWords)
}
