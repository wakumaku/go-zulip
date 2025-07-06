// Package recipient is a package that provides types and functions for handling
// recipients in Zulip messages. It defines interfaces and concrete types for
// representing both channel and direct message recipients.
package recipient

// Recipient For channel messages, either the name or integer ID of the
// channel. For direct messages, either a list containing integer user IDs or
// a list containing string Zulip API email addresses.
type Recipient interface {
	To() any
}

type Channel interface {
	Recipient
	isChannel()
}

type Direct interface {
	Recipient
	isDirect()
}

func ToChannel[T int | string](channel T) Channel {
	switch v := any(channel).(type) {
	case int:
		return toChannelID(v)
	case string:
		return toChannelName(v)
	default:
		return nil
	}
}

type toChannelID int

func (t toChannelID) To() any {
	return t
}

func (t toChannelID) isChannel() {
}

type toChannelName string

func (t toChannelName) To() any {
	return t
}

func (t toChannelName) isChannel() {
}

func ToUser[T int | string](user T) Direct {
	switch v := any(user).(type) {
	case int:
		return toUserIDs([]int{v})
	case string:
		return toUserNames([]string{v})
	default:
		return nil
	}
}

func ToUsers[T []int | []string](users T) Direct {
	switch v := any(users).(type) {
	case []int:
		return toUserIDs(v)
	case []string:
		return toUserNames(v)
	default:
		return nil
	}
}

type toUserIDs []int

func (t toUserIDs) To() any {
	return t
}

func (t toUserIDs) isDirect() {
}

type toUserNames []string

func (t toUserNames) To() any {
	return t
}

func (t toUserNames) isDirect() {
}
