package zulip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNarrower(t *testing.T) {
	cases := []struct {
		name     string
		expected string
		narrower Narrower
	}{
		{
			name:     "empty",
			expected: "[]",
			narrower: Narrower{},
		},
		{
			name:     "is:unread",
			expected: `[{"operator":"is","operand":"unread"}]`,
			narrower: Narrower{}.Add(Is, Unread),
		},
		{
			name:     "is:unread - Operator with operand",
			expected: `[{"operator":"is","operand":"unread"}]`,
			narrower: Narrower{}.Add(IsUnread, nil),
		},
		{
			name:     "is:unread, is:followed",
			expected: `[{"operator":"is","operand":"unread"},{"operator":"is","operand":"followed"}]`,
			narrower: Narrower{}.Add(IsUnread, nil).Add(IsFollowed, nil),
		},
		{
			name:     "channel:1, near:2",
			expected: `[{"operator":"channel","operand":1},{"operator":"near","operand":2}]`,
			narrower: Narrower{}.Add(Channel, 1).Add(Near, 2),
		},
		{
			name:     "channel:1, near:2, not is:unread",
			expected: `[{"operator":"channel","operand":1},{"operator":"near","operand":2},{"operator":"is","operand":"unread","negated":true}]`,
			narrower: Narrower{}.Add(Channel, 1).Add(Near, 2).AddNegated(IsUnread, nil),
		},
		{
			name:     "search:foo, in a dm with user 1",
			expected: `[{"operator":"search","operand":"foo"},{"operator":"dm","operand":1}]`,
			narrower: Narrower{}.Add(Search, "foo").Add(Dm, 1),
		},
		{
			name:     "search:foo, in channel: 1,2,3,4,5",
			expected: `[{"operator":"search","operand":"foo"},{"operator":"channels","operand":[1,2,3,4,5]}]`,
			narrower: Narrower{}.Add(Search, "foo").Add(Channels, []int{1, 2, 3, 4, 5}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.narrower.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expected, string(actual))
		})
	}
}
