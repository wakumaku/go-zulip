package narrow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNarrower(t *testing.T) {
	cases := []struct {
		name          string
		expected      string
		expectedEvent string
		narrower      Filter
	}{
		{
			name:          "empty",
			expected:      "[]",
			expectedEvent: "[]",
			narrower:      NewFilter(),
		},
		{
			name:          "is:unread",
			expected:      `[{"operator":"is","operand":"unread","negated":false}]`,
			expectedEvent: `[["is","unread"]]`,
			narrower:      NewFilter().Add(New(Is, Unread)),
		},
		{
			name:          "is:unread - Operator with operand",
			expected:      `[{"operator":"is","operand":"unread","negated":false}]`,
			expectedEvent: `[["is","unread"]]`,
			narrower:      NewFilter().Add(IsUnread),
		},
		{
			name:          "is:unread, is:followed",
			expected:      `[{"operator":"is","operand":"unread","negated":false},{"operator":"is","operand":"followed","negated":false}]`,
			expectedEvent: `[["is","unread"],["is","followed"]]`,
			narrower:      NewFilter().Add(IsUnread).Add(IsFollowed),
		},
		{
			name:          "channel:1, near:2",
			expected:      `[{"operator":"channel","operand":1,"negated":false},{"operator":"near","operand":2,"negated":false}]`,
			expectedEvent: `[["channel","1"],["near","2"]]`,
			narrower:      NewFilter().Add(New(Channel, 1)).Add(New(Near, 2)),
		},
		{
			name:     "channel:1, near:2, not is:unread",
			expected: `[{"operator":"channel","operand":1,"negated":false},{"operator":"near","operand":2,"negated":false},{"operator":"is","operand":"unread","negated":true}]`,
			narrower: NewFilter().Add(New(Channel, 1)).Add(New(Near, 2)).Add(Negate(IsUnread)),
		},
		{
			name:     "search:foo, in a dm with user 1",
			expected: `[{"operator":"search","operand":"foo","negated":false},{"operator":"dm","operand":1,"negated":false}]`,
			narrower: NewFilter().Add(New(Search, "foo")).Add(New(Dm, 1)),
		},
		{
			name:     "search:foo, in channel: 1,2,3,4,5",
			expected: `[{"operator":"search","operand":"foo","negated":false},{"operator":"channels","operand":[1,2,3,4,5],"negated":false}]`,
			narrower: NewFilter().Add(New(Search, "foo")).Add(New(Channels, []int{1, 2, 3, 4, 5})),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.narrower.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expected, string(actual))

			if tc.expectedEvent != "" {
				eventJSON, err := tc.narrower.MarshalEvent()
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedEvent, string(eventJSON))
			}
		})
	}
}

func TestNewNarrowFromString(t *testing.T) {
	cases := []struct {
		name     string
		narrow   string
		expected Narrow
	}{
		{
			name:     "is:read",
			narrow:   "is:read",
			expected: IsRead,
		},
		{
			name:     "is:read",
			narrow:   "-is:read",
			expected: Negate(IsRead),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			narrow := NewFromString(c.narrow)
			assert.Equal(t, c.expected, narrow)
		})
	}
}

func TestNarrowerBasic(t *testing.T) {
	expectedNarrowerString := "-is:read"

	nf := NewFilter().
		Add(NewNegated("is", "read"))

	assert.Equal(t, expectedNarrowerString, nf.String())

	expectedNarrowersString := "-is:read is:starred"
	expectedNarrowersJSON := `[{"operator":"is", "operand": "read", "negated": true}, {"operator":"is", "operand": "starred", "negated": false}]`
	nfs := NewFilter().
		Add(NewNegated("is", "read")).
		Add(New("is", "starred"))

	assert.Equal(t, expectedNarrowersString, nfs.String())

	currentNarrowersJSON, err := nfs.MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, expectedNarrowersJSON, string(currentNarrowersJSON))
}
