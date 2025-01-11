package main

import (
	"context"
	"log"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/messages"
)

func main() {
	zuliprc, err := zulip.ParseZuliprc(".zuliprc")
	if err != nil {
		log.Fatalf("Error parsing zuliprc file: %v", err)
	}

	apiSection, ok := zuliprc["api"]
	if !ok {
		log.Fatalf("No 'api' section found in zuliprc file")
	}
	log.Printf("Email: %s", apiSection.Email)
	log.Printf("API Key: %s", apiSection.Key)
	log.Printf("Site: %s", apiSection.Site)

	client, err := zulip.NewClient(apiSection.Site, apiSection.Email, apiSection.Key)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	msgs := messages.NewService(client)
	resp, err := msgs.SendMessage(context.TODO(), messages.ToChannelName("general"), "Hello from go-zulip!", messages.ToTopic("test"))
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
	if resp.IsError() {
		log.Fatalf("Error response: %v", resp)
	}

	log.Print("Message sent!")
}
