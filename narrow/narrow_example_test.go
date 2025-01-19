package narrow_test

import (
	"fmt"

	"github.com/wakumaku/go-zulip/narrow"
)

func ExampleFilter() {
	// Create a new filter with a search, topic, and stream narrow
	narrow := narrow.NewFilter().
		Add(narrow.New(narrow.Channel, "general")).
		Add(narrow.New(narrow.Topic, "greetings")).
		Add(narrow.New(narrow.Search, "hello narrow"))

	narrowString := narrow.String()

	fmt.Println(narrowString)
	// Output: channel:general topic:greetings search:hello narrow
}
