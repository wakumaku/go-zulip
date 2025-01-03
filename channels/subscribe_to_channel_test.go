package channels_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/channels"
)

func TestSubscribeToChannel(t *testing.T) {
	client := createMockClient(`{
    "already_subscribed": {
        "1": [
            "testing-help"
        ]
    },
    "msg": "",
    "result": "success",
    "subscribed": {
        "2": [
            "testing-help"
        ]
    }
}`)

	channelSvc := channels.NewService(client)

	msg := map[string]any{
		"subscriptions": `[{"name":"testing-help","description":"test channel"}]`,
	}

	resp, err := channelSvc.SubscribeToChannel(context.Background(),
		[]channels.SubscribeTo{
			{
				Name:        "testing-help",
				Description: "test channel",
			},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Result())

	assert.Equal(t, "testing-help", resp.AlreadySubscribed["1"][0])
	assert.Equal(t, "testing-help", resp.Subscribed["2"][0])

	// validate the parameters sent are correct
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
