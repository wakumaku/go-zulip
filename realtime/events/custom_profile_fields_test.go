package events_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func TestCustomProfileFields(t *testing.T) {
	eventExample := `{
    "fields": [
        {
            "editable_by_user": true,
            "field_data": "",
            "hint": "",
            "id": 1,
            "name": "Phone number",
            "order": 1,
            "required": true,
            "type": 1
        },
        {
            "editable_by_user": true,
            "field_data": "",
            "hint": "What are you known for?",
            "id": 2,
            "name": "Biography",
            "order": 2,
            "required": true,
            "type": 2
        },
        {
            "editable_by_user": true,
            "field_data": "",
            "hint": "Or drink, if you'd prefer",
            "id": 3,
            "name": "Favorite food",
            "order": 3,
            "required": false,
            "type": 1
        },
        {
            "display_in_profile_summary": true,
            "editable_by_user": true,
            "field_data": "{\"0\":{\"text\":\"Vim\",\"order\":\"1\"},\"1\":{\"text\":\"Emacs\",\"order\":\"2\"}}",
            "hint": "",
            "id": 4,
            "name": "Favorite editor",
            "order": 4,
            "required": true,
            "type": 3
        },
        {
            "editable_by_user": false,
            "field_data": "",
            "hint": "",
            "id": 5,
            "name": "Birthday",
            "order": 5,
            "required": false,
            "type": 4
        },
        {
            "display_in_profile_summary": true,
            "editable_by_user": true,
            "field_data": "",
            "hint": "Or your personal blog's URL",
            "id": 6,
            "name": "Favorite website",
            "order": 6,
            "required": false,
            "type": 5
        },
        {
            "editable_by_user": false,
            "field_data": "",
            "hint": "",
            "id": 7,
            "name": "Mentor",
            "order": 7,
            "required": true,
            "type": 6
        },
        {
            "editable_by_user": true,
            "field_data": "{\"subtype\":\"github\"}",
            "hint": "Enter your GitHub username",
            "id": 8,
            "name": "GitHub",
            "order": 8,
            "required": true,
            "type": 7
        }
    ],
    "id": 0,
    "type": "custom_profile_fields"
}`

	v := events.CustomProfileFields{}
	err := json.Unmarshal([]byte(eventExample), &v)
	assert.NoError(t, err)

	assert.Equal(t, 0, v.EventID())
	assert.Equal(t, events.CustomProfileFieldsType, v.EventType())
	assert.Equal(t, "custom_profile_fields", v.EventOp())

	assert.Len(t, v.Fields, 8)
	assert.True(t, v.Fields[0].EditableByUser)
	assert.Equal(t, "", v.Fields[0].FieldData)
	assert.Equal(t, "", v.Fields[0].Hint)
	assert.Equal(t, 1, v.Fields[0].ID)
	assert.Equal(t, "Phone number", v.Fields[0].Name)
	assert.Equal(t, 1, v.Fields[0].Order)
	assert.True(t, v.Fields[0].Required)
	assert.Equal(t, 1, v.Fields[0].Type)
}
