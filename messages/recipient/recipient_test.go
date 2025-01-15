package recipient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipients(t *testing.T) {
	recipientChannelString := ToChannel("channel")
	assert.Equal(t, toChannelName("channel"), recipientChannelString.To())

	recipientChannelInt := ToChannel(42)
	assert.Equal(t, toChannelID(42), recipientChannelInt.To())

	recipientDirectString := ToUser("john")
	assert.Equal(t, toUserNames([]string{"john"}), recipientDirectString.To())

	recipientDirectInt := ToUser(42)
	assert.Equal(t, toUserIDs([]int{42}), recipientDirectInt.To())

	recipientDirectStrings := ToUsers([]string{"john", "doe"})
	assert.Equal(t, toUserNames([]string{"john", "doe"}), recipientDirectStrings.To())

	recipientDirectInts := ToUsers([]int{42, 43})
	assert.Equal(t, toUserIDs([]int{42, 43}), recipientDirectInts.To())
}
