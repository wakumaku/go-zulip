package invitations_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip/invitations"
)

func TestCreateReusableInvitationLink(t *testing.T) {
	client := createMockClient(`{
    "invite_link": "https://example.zulipchat.com/join/yddhtzk4jgl7rsmazc5fyyyy/",
    "msg": "",
    "result": "success"
}`)

	service := invitations.NewService(client)

	resp, err := service.CreateReusableInvitationLink(context.Background(), invitations.StreamIds([]int{1, 2}))
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Result())
	assert.Equal(t, "https://example.zulipchat.com/join/yddhtzk4jgl7rsmazc5fyyyy/", resp.InviteLink)

	// validate the parameters sent are correct
	assert.Equal(t, map[string]any{
		"stream_ids": "[1,2]",
	}, client.(*mockClient).paramsSent)
}
