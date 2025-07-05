package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/wakumaku/go-zulip/channels"
	"github.com/wakumaku/go-zulip/messages"
	"github.com/wakumaku/go-zulip/messages/recipient"
	"github.com/wakumaku/go-zulip/narrow"
	"github.com/wakumaku/go-zulip/realtime"
	"github.com/wakumaku/go-zulip/realtime/events"
	"github.com/wakumaku/go-zulip/users"
)

// Bot holds the services to interact with Zulip's API
type Bot struct {
	UserSVC     *users.Service
	ChannelSVC  *channels.Service
	MessageSVC  *messages.Service
	RealtimeSVC *realtime.Service
}

func (b *Bot) Run(ctx context.Context, channel, topic string) error {
	// Subscribe to the channel
	if subscriptionResponse, err := b.ChannelSVC.SubscribeToChannel(ctx,
		[]channels.SubscribeTo{{Name: channel}},
	); err != nil || subscriptionResponse.IsError() {
		if err != nil {
			return fmt.Errorf("failed to subscribe to channel: %v", err)
		}

		return fmt.Errorf("failed to subscribe to channel: %v", subscriptionResponse.Msg())
	}

	// Identify itself
	respUserMe, err := b.UserSVC.GetUserMe(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user me: %v", err)
	}

	if respUserMe.IsError() {
		return fmt.Errorf("zulip API error getting user me: %v", respUserMe.Msg())
	}

	botID := respUserMe.UserID
	botEmail := respUserMe.Email
	botName := respUserMe.FullName

	// Register a queue for the bot, will receive only messages from the subscribed channel
	queueRegisterResp, err := b.RealtimeSVC.RegisterEvetQueue(ctx,
		realtime.EventTypes(events.MessageType),
		realtime.NarrowEvents(narrow.NewFilter().
			Add(narrow.New(narrow.Stream, channel)),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to register event queue: %v", err)
	}

	if queueRegisterResp.IsError() {
		return fmt.Errorf("zulip API error registering event queue: %v", queueRegisterResp.Msg())
	}

	queueID := queueRegisterResp.QueueId
	lastMessageID := queueRegisterResp.LastEventId

	log.Printf("Bot ID: %d - %s (%s) is ready", botID, botName, botEmail)

	messages := b.messageEvents(ctx, queueID, lastMessageID)

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				return ctx.Err()
			}

			return nil

		case <-time.After(10 * time.Second):
			// Sends a nice good morning image after 10 seconds waiting for someone to say something ...
			uploadedFile, err := b.MessageSVC.UploadFile(ctx, "./good_morning.png")
			if err != nil {
				return fmt.Errorf("failed to upload file: %v", err)
			}

			if uploadedFile.IsError() {
				return fmt.Errorf("zulip API error uploading file: %v", uploadedFile.Msg())
			}

			helloThere := fmt.Sprintf("Hello there! Am I alone in #**%s**? My name is @**%s**",
				channel, botName)

			sendMsgResp, err := b.MessageSVC.SendMessageToChannelTopic(ctx, recipient.ToChannel(channel), topic, helloThere)
			if err != nil {
				return fmt.Errorf("failed to send message: %v", err)
			}

			if sendMsgResp.IsError() {
				return fmt.Errorf("zulip API error sending message: %v", sendMsgResp.Msg())
			}

			// Sends a nice good morning image
			sendImgResp, err := b.MessageSVC.SendMessageToChannelTopic(ctx, recipient.ToChannel(channel), topic, fmt.Sprintf("[good_morning.png](%s)", uploadedFile.URI))
			if err != nil {
				return fmt.Errorf("failed to send image message: %v", err)
			}

			if sendImgResp.IsError() {
				return fmt.Errorf("zulip API error sending image message: %v", sendImgResp.Msg())
			}

			log.Printf("Bot ID: %d - %s (%s) - Started the conversation", botID, botName, botEmail)

		case e, ok := <-messages:
			// Manage received messages to the channel
			if !ok {
				// log.Printf("Bot ID: %d - Messages channel closed", botID)
				return nil
			}

			if e.Message.SenderID == botID {
				// log.Printf("Bot ID: %d - Ignoring self message", botID)
				continue
			}

			// Simulate typing ... :)
			time.Sleep(2 * time.Second)

			responseMessage, reaction := b.makeResponse(e.Message.SenderFullName, e.Message.Content)
			if reaction != "" {
				// Add reaction to the message
				reactionResp, err := b.MessageSVC.AddEmojiReaction(ctx, e.Message.ID, reaction)
				if err != nil {
					return fmt.Errorf("failed to add reaction: %v", err)
				}

				if reactionResp.IsError() {
					return fmt.Errorf("zulip API error adding reaction: %v", reactionResp.Msg())
				}

				log.Printf("Bot ID: %d - %s (%s) - Reacted to a message", botID, botName, botEmail)
			}

			if responseMessage == "" {
				continue
			}

			if responseMessage == "SLEEP" {
				time.Sleep(5 * time.Second)
				continue
			}

			sendMsgResp, err := b.MessageSVC.SendMessageToChannelTopic(ctx, recipient.ToChannel(channel), topic, responseMessage)
			if err != nil {
				return fmt.Errorf("failed to send message: %v", err)
			}

			if sendMsgResp.IsError() {
				return fmt.Errorf("zulip API error sending message: %v", sendMsgResp.Msg())
			}

			log.Printf("Bot ID: %d - %s (%s) - Replied to a message", botID, botName, botEmail)
		}
	}
}

// makeResponse returns the response message and the reaction to be added to the message depending on the received
// message in a very simple way.
// If the message received starts or ends with a specific text, the response message and the reaction are returned.
func (b *Bot) makeResponse(from, message string) (string, string) {
	if strings.HasPrefix(message, "Hello there! Am I alone") {
		return fmt.Sprintf("Of course not! I'm here with you, @**%s**", from), "wave"
	}

	if strings.HasPrefix(message, "Of course not! I'm here with you") {
		return "Aaaah, glad to see you here", ""
	}

	if strings.HasSuffix(message, "Aaaah, glad to see you here") {
		return "Cool! And now what ... :question:", ""
	}

	if strings.HasSuffix(message, "Cool! And now what ... :question:") {
		return fmt.Sprintf("@**%s** I dont't know, I guess this is enough. I will go to sleep :sleeping:", from), "thinking"
	}

	if strings.HasSuffix(message, "I will go to sleep :sleeping:") {
		return "Ooook, me too. See ya!", ""
	}

	if strings.HasSuffix(message, "Ooook, me too. See ya!") {
		return "SLEEP", "goodbye"
	}

	return "", ""
}

// messageEvents returns a channel with the messages received from the event queue
// If something goes wrong, the error is logged and the channel is closed.
func (b *Bot) messageEvents(ctx context.Context, queueID string, lastMessageID int) chan *events.Message {
	messages := make(chan *events.Message)

	go func() {
		defer close(messages)

		for {
			eventsResp, err := b.RealtimeSVC.GetEventsEventQueue(ctx, queueID, realtime.LastEventID(lastMessageID))
			if err != nil {
				log.Printf("failed to get events from event queue: %v", err)
				return
			}

			if eventsResp.IsError() {
				log.Printf("zulip API error getting events from event queue: %v", eventsResp.Msg())
				return
			}

			for _, event := range eventsResp.Events {
				lastMessageID = event.EventID()

				switch e := event.(type) {
				case *events.Message:
					messages <- e
				default:
					continue
				}
			}
		}
	}()

	return messages
}
