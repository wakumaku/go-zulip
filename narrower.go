package zulip

import (
	"encoding/json"
	"fmt"
)

// NarrowExpression is an interface that represents the operators that can be used in the narrow parameter.
type NarrowExpression interface {
	// withOperand returns true if the operator is a composite operator in form of "is:unread".
	withOperand() bool
}

// NarrowOperator is a type that represents the concrete operators that can be used in the narrow parameter.
type NarrowOperator string

func (n NarrowOperator) withOperand() bool {
	return false
}

func (n NarrowOperator) String() string {
	return string(n)
}

// NarrowOperand is a type that represents the concrete operands that can be used in the narrow parameter.
type NarrowOperand any

type NarrowOperatorOperand struct {
	operator NarrowOperator
	operand  NarrowOperand
}

func (n NarrowOperatorOperand) withOperand() bool {
	return true
}

func (n NarrowOperatorOperand) String() string {
	return fmt.Sprintf("%s:%s", n.operator, n.operand)
}

// Operators List
const (
	// Search for only the message with ID 12345
	Id NarrowOperator = "id"
	// Search for the conversation that contains the message with ID 12345.
	With     NarrowOperator = "with"
	Near     NarrowOperator = "near"
	Channel  NarrowOperator = "channel"
	Channels NarrowOperator = "channels"
	Stream   NarrowOperator = "stream"
	Streams  NarrowOperator = "streams"
	Sender   NarrowOperator = "sender"
	Search   NarrowOperator = "search"
	// Search the direct message conversation between you and user ID 1234
	Dm NarrowOperator = "dm"
	// Search all direct messages (1-on-1 and group) that include you and user ID 1234.
	DmIncluding NarrowOperator = "dm-including"

	Is  NarrowOperator = "is"
	Has NarrowOperator = "has"
)

// Operands List
var (
	Unread   NarrowOperand = NarrowOperand("unread")
	Followed NarrowOperand = NarrowOperand("followed")
	// Dm         NarrowOperand = NarrowOperand("dm") // clashes with Dm NarrowOperator
	Mentioned  NarrowOperand = NarrowOperand("mentioned")
	Starred    NarrowOperand = NarrowOperand("starred")
	Read       NarrowOperand = NarrowOperand("read")
	Alerted    NarrowOperand = NarrowOperand("alerted")
	Attachment NarrowOperand = NarrowOperand("attachment")
	Image      NarrowOperand = NarrowOperand("image")
	Link       NarrowOperand = NarrowOperand("link")
	Reaction   NarrowOperand = NarrowOperand("reaction")
)

// Operators with Operands
var (
	IsUnread      NarrowOperatorOperand = NarrowOperatorOperand{Is, Unread}
	IsFollowed    NarrowOperatorOperand = NarrowOperatorOperand{Is, Followed}
	IsDm          NarrowOperatorOperand = NarrowOperatorOperand{Is, NarrowOperand("dm")}
	IsMentioned   NarrowOperatorOperand = NarrowOperatorOperand{Is, Mentioned}
	IsStarred     NarrowOperatorOperand = NarrowOperatorOperand{Is, Starred}
	IsRead        NarrowOperatorOperand = NarrowOperatorOperand{Is, Read}
	IsAlerted     NarrowOperatorOperand = NarrowOperatorOperand{Is, Alerted}
	HasAttachment NarrowOperatorOperand = NarrowOperatorOperand{Is, Attachment}
	HasImage      NarrowOperatorOperand = NarrowOperatorOperand{Is, Image}
	HasLink       NarrowOperatorOperand = NarrowOperatorOperand{Has, Link}
	HasReaction   NarrowOperatorOperand = NarrowOperatorOperand{Has, Reaction}
)

// Narrower contains a list of constraints to narrow down the search results.
type Narrower []NarrowItem

// NarrowItem represents a single constraint in the Narrower.
type NarrowItem struct {
	Operator NarrowExpression `json:"operator"`
	Operand  NarrowOperand    `json:"operand,omitempty"`
	Negated  *bool            `json:"negated,omitempty"`
}

// Add adds a new constraint to the Narrower.
// Operator can be a composite operator like "is:unread", therefore operands will be ignored.
func (n Narrower) Add(operator NarrowExpression, operand NarrowOperand) Narrower {
	if operator.withOperand() {
		operator, operand = getOperatorOperand(operator)
	}

	return append(n, NarrowItem{Operator: operator, Operand: operand})
}

// AddNegated adds a new constraint to the Narrower with negation.
func (n Narrower) AddNegated(operator NarrowExpression, operand NarrowOperand) Narrower {
	if operator.withOperand() {
		operator, operand = getOperatorOperand(operator)
	}

	negated := true
	return append(n, NarrowItem{Operator: operator, Operand: operand, Negated: &negated})
}

// MarshalJSON implements the json.Marshaler interface which is used to perform Message search.
// ref: https://zulip.com/api/construct-narrow
func (n Narrower) MarshalJSON() ([]byte, error) {
	return json.Marshal([]NarrowItem(n))
}

// EventsJSON returns the JSON representation of the Narrower for the register-queue endpoint.
// ref: https://zulip.com/api/register-queue#parameter-narrow
func (n Narrower) EventsJSON() ([]byte, error) {
	out := make([][]string, 0, len(n))
	for _, item := range n {
		operator, operand := item.Operator, item.Operand
		if item.Operator.withOperand() {
			operator, operand = getOperatorOperand(item.Operator)
		}
		out = append(out, []string{fmt.Sprintf("%s", operator), fmt.Sprintf("%s", operand)})
	}
	return json.Marshal(out)
}

// getOperatorOperand returns the operator and operand from the NarrowOperatorOperand.
func getOperatorOperand(n NarrowExpression) (NarrowExpression, any) {
	op, ok := n.(NarrowOperatorOperand)
	if !ok {
		return n, nil
	}

	return op.operator, op.operand
}
