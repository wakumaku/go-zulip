package messages_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/messages"
	"github.com/wakumaku/go-zulip/narrow"
)

func TestUpdatePersonalMessageFlagsNarrow(t *testing.T) {
	client := createMockClient(`{
    "first_processed_id": 35,
    "found_newest": true,
    "found_oldest": false,
    "last_processed_id": 55,
    "msg": "",
    "processed_count": 11,
    "result": "success",
    "updated_count": 8
}`)

	msgSvc := messages.NewService(client)

	resp, err := msgSvc.UpdatePersonalMessageFlagsNarrow(
		context.Background(),
		"anchor",
		10,
		10,
		narrow.NewFilter().Add(narrow.New(narrow.Channel, "Denmark")),
		messages.OperationAdd,
		messages.FlagRead,
		messages.UpdatePersonalMessageFlagsNarrowIncludeAnchor(),
	)
	assert.NoError(t, err)
	assert.Equal(t, 35, resp.FirstProcessedID)
	assert.Equal(t, 55, resp.LastProcessedID)
	assert.Equal(t, 11, resp.ProcessedCount)
	assert.Equal(t, 8, resp.UpdatedCount)
	assert.True(t, resp.FoundNewest)
	assert.False(t, resp.FoundOldest)

	// validate payload
	expedtedParams := map[string]interface{}{
		"anchor":         "anchor",
		"num_before":     10,
		"num_after":      10,
		"narrow":         `[{"operator":"channel","operand":"Denmark","negated":false}]`,
		"op":             messages.OperationAdd,
		"flag":           messages.FlagRead,
		"include_anchor": true,
	}
	assert.Equal(t, expedtedParams, client.(*mockClient).paramsSent)
}
