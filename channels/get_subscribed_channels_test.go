package channels_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/channels"
)

func TestGetSubscribedChannels(t *testing.T) {
	client := createMockClient(`{
    "msg": "",
    "result": "success",
    "subscriptions": [
        {
            "audible_notifications": true,
            "color": "#e79ab5",
            "creator_id": null,
            "description": "A Scandinavian country",
            "desktop_notifications": true,
            "invite_only": false,
            "is_archived": false,
            "is_muted": false,
            "name": "Denmark",
            "pin_to_top": false,
            "push_notifications": false,
            "stream_id": 1,
            "subscribers": [
                7,
                10,
                11,
                12,
                14
            ]
        },
        {
            "audible_notifications": true,
            "color": "#e79ab5",
            "creator_id": 8,
            "description": "Located in the United Kingdom",
            "desktop_notifications": true,
            "invite_only": false,
            "is_archived": false,
            "is_muted": false,
            "name": "Scotland",
            "pin_to_top": false,
            "push_notifications": false,
            "stream_id": 3,
            "subscribers": [
                7,
                11,
                12,
                14
            ]
        }
    ]
}`)

	channelSvc := channels.NewService(client)

	msg := map[string]any{
		"include_subscribers": true,
	}

	resp, err := channelSvc.GetSubscribedChannels(context.Background(),
		channels.IncludeSubscribersList(true),
	)
	require.NoError(t, err)
	assert.Len(t, resp.Subscriptions, 2)
	assert.Equal(t, "Denmark", resp.Subscriptions[0].Name)
	assert.Equal(t, "Scotland", resp.Subscriptions[1].Name)
	assert.Equal(t, []int{7, 10, 11, 12, 14}, resp.Subscriptions[0].Subscribers)

	// validate the parameters sent are correct
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
