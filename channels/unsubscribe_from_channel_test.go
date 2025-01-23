package channels_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/channels"
)

func TestUnsubscribeFromChannel(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "not_removed": [],
    "removed": [
        "testing-help"
    ],
    "result": "success"
}`)

	channelSvc := channels.NewService(client)

	resp, err := channelSvc.UnsubscribeFromChannel(context.Background(), []string{"testing-help"}, channels.Principals([]int{1, 2}))
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Result())

	// validate the parameters sent are correct
	assert.Equal(t, map[string]any{
		"subscriptions": "[\"testing-help\"]",
		"principals":    "[1,2]",
	}, client.(*mockClient).paramsSent)
}
