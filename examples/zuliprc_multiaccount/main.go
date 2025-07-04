package main

import (
	"context"
	"log"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/messages"
	"github.com/wakumaku/go-zulip/messages/recipient"
)

func main() {
	// Read the zuliprc file on specific section for each bot
	credentialsBot1 := zulip.CredentialsFromZuliprc("zuliprc", "bot1")
	credentialsBot2 := zulip.CredentialsFromZuliprc("zuliprc", "bot2")
	credentialsBot3 := zulip.CredentialsFromZuliprc("zuliprc", "bot3")

	// Create a client for each bot
	bot1, err := zulip.NewClient(credentialsBot1)
	if err != nil {
		log.Fatalf("Error creating bot1: %v", err)
	}

	bot2, err := zulip.NewClient(credentialsBot2)
	if err != nil {
		log.Fatalf("Error creating bot2: %v", err)
	}

	bot3, err := zulip.NewClient(credentialsBot3)
	if err != nil {
		log.Fatalf("Error creating bot3: %v", err)
	}

	// Create a message service for each bot
	msgBot1 := messages.NewService(bot1)
	msgBot2 := messages.NewService(bot2)
	msgBot3 := messages.NewService(bot3)

	ctx := context.Background()

	// Send a message to a channel/topic for each bot
	resp, err := msgBot1.SendMessageToChannelTopic(ctx, recipient.ToChannel("general"), "greetings", "I'm bot 1!")
	if err != nil {
		log.Fatalf("[bot1] Error sending message: %v", err)
	}
	if resp.IsError() {
		log.Fatalf("[bot1] Error response: %v", resp)
	}

	resp, err = msgBot2.SendMessageToChannelTopic(ctx, recipient.ToChannel("general"), "greetings", "I'm bot 2!")
	if err != nil {
		log.Fatalf("[bot2] Error sending message: %v", err)
	}
	if resp.IsError() {
		log.Fatalf("[bot2] Error response: %v", resp)
	}

	resp, err = msgBot3.SendMessageToChannelTopic(ctx, recipient.ToChannel("general"), "greetings", "I'm bot 3!")
	if err != nil {
		log.Fatalf("[bot3] Error sending message: %v", err)
	}
	if resp.IsError() {
		log.Fatalf("[bot3] Error response: %v", resp)
	}

	log.Print("Messages sent!")
}
