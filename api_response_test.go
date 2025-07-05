package zulip

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIResponseMarshalling(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{
			name: "Normal response",
			input: `{
    "code": "INVALID_API_KEY",
    "msg": "Invalid API key",
    "result": "error"
}`,
		},
		{
			name: "Extra fields (1)",
			input: `{
    "code": "REQUEST_VARIABLE_MISSING",
    "msg": "Missing 'content' argument",
    "result": "error",
    "var_name": "content"
}`,
		},
		{
			name: "Extra fields (2)",
			input: `{
    "ignored_parameters_unsupported": [
        "invalid_param_1",
        "invalid_param_2"
    ],
    "msg": "",
    "result": "success"
}`,
		},
	}

	for _, c := range cases {
		er := APIResponseBase{}
		require.NoError(t, json.Unmarshal([]byte(c.input), &er))

		erJSON, err := json.Marshal(er)
		require.NoError(t, err)
		assert.JSONEq(t, c.input, string(erJSON))
	}
}
