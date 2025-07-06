// Package narrow provides functionality for creating message search filters in Zulip.
// This includes operators and operands for filtering messages by various criteria
// such as sender, channel, topic, message content, and message properties.
// Narrow filters are used in message queries and event queue registration.
//
// See https://zulip.com/api/ for the complete API documentation.
package narrow

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Operator is a type that represents the concrete operators that can be used in the narrow parameter.
type Operator string

func (n Operator) String() string {
	return string(n)
}

// Operand is a type that represents the concrete operands that can be used in the narrow parameter.
type Operand any

// Operators List
const (
	// Search for only the message with ID 12345
	ID Operator = "id"
	// Search for the conversation that contains the message with ID 12345.
	With     Operator = "with"
	Near     Operator = "near"
	Channel  Operator = "channel"
	Channels Operator = "channels"
	// stream is a legacy alias for "channel"
	Stream Operator = "stream"
	// streams is a legacy alias for "channels"
	Streams Operator = "streams"
	Topic   Operator = "topic"
	Sender  Operator = "sender"
	Search  Operator = "search"
	// Search the direct message conversation between you and user ID 1234
	Dm Operator = "dm"
	// Search all direct messages (1-on-1 and group) that include you and user ID 1234.
	DmIncluding Operator = "dm-including"

	Is  Operator = "is"
	Has Operator = "has"
)

// Operands List
var (
	Unread   Operand = "unread"
	Followed Operand = "followed"
	// Dm         Operand = "dm" // clashes with Dm Operator
	Mentioned  Operand = "mentioned"
	Starred    Operand = "starred"
	Read       Operand = "read"
	Alerted    Operand = "alerted"
	Attachment Operand = "attachment"
	Image      Operand = "image"
	Link       Operand = "link"
	Reaction   Operand = "reaction"
)

// Operators with Operands
var (
	IsUnread      Narrow = New(Is, Unread)
	IsFollowed    Narrow = New(Is, Followed)
	IsDm          Narrow = New(Is, Operand("dm"))
	IsMentioned   Narrow = New(Is, Mentioned)
	IsStarred     Narrow = New(Is, Starred)
	IsRead        Narrow = New(Is, Read)
	IsAlerted     Narrow = New(Is, Alerted)
	HasAttachment Narrow = New(Is, Attachment)
	HasImage      Narrow = New(Is, Image)
	HasLink       Narrow = New(Has, Link)
	HasReaction   Narrow = New(Has, Reaction)
)

// Filter is a collection of narrow items to be applied when searching for
// messages or filtering events
type Filter []Narrow

// NewFilter creates a new Filter
func NewFilter() Filter {
	return make(Filter, 0)
}

// Add adds a new narrow item to the filter and returns the updated filter
func (f Filter) Add(narrow Narrow) Filter {
	return append(f, narrow)
}

// String returns a string representation of the Filter
func (f *Filter) String() string {
	ns := make([]string, len(*f))
	for i, n := range *f {
		ns[i] = n.String()
	}

	return strings.Join(ns, " ")
}

// MarshalJSON returns the JSON encoding of the Filter
func (f Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal([]Narrow(f))
}

// MarshalEvent returns the JSON encoding of the Filter for events
func (f Filter) MarshalEvent() ([]byte, error) {
	out := make([][]string, 0, len(f))
	for _, item := range f {
		operator, operand := item.Operator, item.Operand
		out = append(out, []string{string(operator), fmt.Sprintf("%v", operand)})
	}

	return json.Marshal(out)
}

// Narrow is a struct that represents a single narrow item to be applied when searching for messages or filtering events
type Narrow struct {
	Operator Operator `json:"operator"`
	Operand  Operand  `json:"operand"`
	Negated  bool     `json:"negated"`
}

// New creates a new Narrow
func New(op Operator, val Operand) Narrow {
	return newNarrow(op, val, false)
}

// NewNegated creates a new Negated Narrow
func NewNegated(op Operator, val Operand) Narrow {
	return newNarrow(op, val, true)
}

// NewFromString creates a new Narrow from a string
func NewFromString(s string) Narrow {
	s = strings.TrimSpace(s)

	isNegated := strings.HasPrefix(s, "-")

	opetatorOperand := strings.TrimPrefix(s, "-")

	operatorOperandSlice := strings.Split(opetatorOperand, ":")
	if len(operatorOperandSlice) != 2 {
		return Narrow{}
	}

	operator, operand := operatorOperandSlice[0], operatorOperandSlice[1]

	if isNegated {
		return NewNegated(Operator(operator), Operand(operand))
	}

	return New(Operator(operator), Operand(operand))
}

func newNarrow(op Operator, val Operand, negated bool) Narrow {
	return Narrow{
		Operator: op,
		Operand:  val,
		Negated:  negated,
	}
}

// Negate returns a new Negated Narrow
func Negate(n Narrow) Narrow {
	n.negate()
	return n
}

func (n *Narrow) negate() {
	n.Negated = true
}

// String returns a string representation of the Narrow
func (n *Narrow) String() string {
	negated := ""
	if n.Negated {
		negated = "-"
	}

	return fmt.Sprintf("%s%s:%s", negated, n.Operator, n.Operand)
}
