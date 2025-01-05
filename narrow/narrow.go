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
	Id Operator = "id"
	// Search for the conversation that contains the message with ID 12345.
	With     Operator = "with"
	Near     Operator = "near"
	Channel  Operator = "channel"
	Channels Operator = "channels"
	Stream   Operator = "stream"
	Streams  Operator = "streams"
	Topic    Operator = "topic"
	Sender   Operator = "sender"
	Search   Operator = "search"
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
	// Dm         NarrowOperand = NarrowOperand("dm") // clashes with Dm NarrowOperator
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

// Narrow filter is a collection of filters to be applied when searching for
// messages or filtering events
type Filter []Narrow

func NewFilter() Filter {
	return make(Filter, 0)
}

func (nf Filter) Add(narrow Narrow) Filter {
	return append(nf, narrow)
}

func (nf *Filter) String() string {
	ns := make([]string, len(*nf))
	for i, n := range *nf {
		ns[i] = n.String()
	}

	return strings.Join(ns, " ")
}

func (n Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal([]Narrow(n))
}

func (n Filter) MarshalEvent() ([]byte, error) {
	out := make([][]string, 0, len(n))
	for _, item := range n {
		operator, operand := item.Operator, item.Operand
		out = append(out, []string{string(operator), fmt.Sprintf("%v", operand)})
	}
	return json.Marshal(out)
}

// Narrow
type Narrow struct {
	Operator Operator `json:"operator"`
	Operand  Operand  `json:"operand"`
	Negated  bool     `json:"negated"`
}

func New(op Operator, val Operand) Narrow {
	return newNarrow(op, val, false)
}

func NewNegated(op Operator, val Operand) Narrow {
	return newNarrow(op, val, true)
}

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

func Negate(n Narrow) Narrow {
	n.negate()
	return n
}

func (n *Narrow) negate() {
	n.Negated = true
}

// Stringer
func (n *Narrow) String() string {
	negated := ""
	if n.Negated {
		negated = "-"
	}

	return fmt.Sprintf("%s%s:%s", negated, n.Operator, n.Operand)
}
