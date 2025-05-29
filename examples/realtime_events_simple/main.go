package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/realtime"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func main() {
	ctx := context.Background()

	// Initialize client
	credentials := zulip.Credentials("https://chat.zulip.org", "email@zulip.org", "0123456789")
	c, err := zulip.NewClient(credentials)
	if err != nil {
		log.Fatal(err)
	}

	realtimeSvc := realtime.NewService(c)

	// Register a queue passing the events we want to receive
	queue, err := realtimeSvc.RegisterEvetQueue(ctx,
		realtime.EventTypes(
			events.AlertWordsType,
			events.AttachmentType,
			events.MessageType,
			events.PresenceType,
			events.RealmEmojiType,
			events.RealmUserType,
			events.SubmessageType,
			events.TypingType,
			events.UpdateMessageType,
			events.DeleteMessageType,
		),
		realtime.AllPublicStreams(true),
	)
	if err != nil {
		log.Fatalf("error registering event queue: %s", err)
	}

	if queue.IsError() {
		log.Fatalf("%s: %s", queue.Msg(), queue.Code())
	}

	log.Printf("QueueId: %s", queue.QueueId)
	log.Println("Waiting for events...")

	lastEventID := queue.LastEventId

	// Infinite loop polling for new events
	for {
		// Long polling HTTP Request
		eventsFromQueue, err := realtimeSvc.GetEventsEventQueue(ctx, queue.QueueId, realtime.LastEventID(lastEventID))
		if err != nil {
			log.Fatalf("error getting events from queue: %s", err)
		}

		for _, e := range eventsFromQueue.Events {
			var logEntry string
			// Identify the message type received
			switch event := e.(type) {
			case *events.Message:
				if event.Message.DisplayRecipient.IsChannel {
					logEntry = fmt.Sprintf("#%s [%s]: %s", event.Message.DisplayRecipient.Channel, event.Message.SenderFullName, event.Message.Content)
				} else {
					var users []string
					for _, user := range event.Message.DisplayRecipient.Users {
						users = append(users, user.FullName)
					}
					logEntry = fmt.Sprintf("@%s: %s", strings.Join(users, ", @"), event.Message.Content)
				}

			case *events.AlertWords:
				logEntry = fmt.Sprintf("!AlertWords ID: %d, Words: %s", event.ID, event.AlertWords)

			case *events.RealmUser:
				logEntry = fmt.Sprintf("@RealmUser ID: %d, Op: %s, FullName: %s", event.ID, event.Op, event.Person.FullName)

			case *events.Presence:
				logEntry = fmt.Sprintf("*Presence Email: %s, Status: %s", event.Email, event.Presence.Website.Status)

			case *events.RealmEmoji:
				logEntry = fmt.Sprintf(":RealmEmoji event ID: %d\n", event.ID)
				for id, emoji := range event.RealmEmoji {
					logEntry += fmt.Sprintf("  %s: %s, %s\n", id, emoji.Name, emoji.SourceURL)
				}

			default:
				logEntry = fmt.Sprintf("#%d %s", e.EventID(), e.EventType())
			}

			log.Println(logEntry)

			lastEventID = e.EventID()
		}
	}
}
