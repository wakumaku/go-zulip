package main

import (
	"context"
	"log"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/messages"
	"github.com/wakumaku/go-zulip/messages/recipient"
)

func main() {
	ctx := context.Background()

	// Initialize client
	credentials := zulip.Credentials("https://chat.zulip.org", "email@zulip.org", "0123456789")
	c, err := zulip.NewClient(credentials)
	if err != nil {
		log.Fatal(err)
	}

	// Send a message to a channel/topic
	msgSvc := messages.NewService(c)
	sendMessageResponse, err := msgSvc.SendMessageToChannelTopic(ctx,
		recipient.ToChannel("general"), "greetings",
		"Hello Zulip!",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the previous message
	fetchMessageResponse, err := msgSvc.FetchSingleMessage(ctx, sendMessageResponse.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fetchMessageResponse.Message.ID, fetchMessageResponse.Message.Content)
}
