package events

const CustomProfileFieldsType EventType = "custom_profile_fields"

type CustomProfileFields struct {
	ID   int       `json:"id"`
	Type EventType `json:"type"`
	CustomProfileFieldsData
}

type CustomProfileFieldsData struct {
	Fields []CustomProfileField `json:"fields"`
}

type CustomProfileField struct {
	EditableByUser bool   `json:"editable_by_user"`
	FieldData      string `json:"field_data"`
	Hint           string `json:"hint"`
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Order          int    `json:"order"`
	Required       bool   `json:"required"`
	Type           int    `json:"type"`
}

func (e *CustomProfileFields) EventID() int {
	return e.ID
}

func (e *CustomProfileFields) EventType() EventType {
	return e.Type
}

func (e *CustomProfileFields) EventOp() string {
	return string(CustomProfileFieldsType)
}
