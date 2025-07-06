package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/messages"
	"github.com/wakumaku/go-zulip/messages/recipient"
	"github.com/wakumaku/go-zulip/realtime"
	"github.com/wakumaku/go-zulip/realtime/events"
	"github.com/wakumaku/go-zulip/users"
)

// Parrot is a bot that repeats messages it receives.
// Every message it can "hear" is sent back to the same recipient
// If the message is sent to a channel, the message is sent back via direct message

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	insecureClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	z, err := zulip.NewClient(zulip.CredentialsFromZuliprc("zuliprc", "api"), zulip.WithHTTPClient(&insecureClient))
	if err != nil {
		panic(err)
	}

	userSvc := users.NewService(z)
	msgSvc := messages.NewService(z)
	rtimeSvc := realtime.NewService(z)

	parrot, err := userSvc.GetUserMe(ctx)
	if err != nil {
		panic(err)
	}

	queue, err := rtimeSvc.RegisterEvetQueue(ctx,
		realtime.EventTypes(events.MessageType),
		realtime.AllPublicStreams(true),
	)
	if err != nil {
		panic(err)
	}

	if queue.IsError() {
		panic(queue.Msg())
	}

	lastEventID := queue.LastEventID
	for {
		evs, err := rtimeSvc.GetEventsEventQueue(ctx, queue.QueueID, realtime.LastEventID(lastEventID))
		if err != nil {
			panic(err)
		}

		if evs.IsError() {
			panic(evs.Msg())
		}

		log.Printf("Received Events: %d\n", len(evs.Events))

		for _, ev := range evs.Events {
			lastEventID = ev.EventID()

			if e, ok := ev.(*events.Message); ok {
				if e.Message.IsMeMessage {
					continue
				}

				if e.Message.SenderID == parrot.UserID {
					continue
				}

				whereYouSaid := fmt.Sprintf("in #%s (%s)", e.Message.DisplayRecipient.Channel, e.Message.Subject)
				if !e.Message.DisplayRecipient.IsChannel {
					whereYouSaid = "to me"
				}

				resp, err := msgSvc.SendMessageToUsers(ctx, recipient.ToUser(e.Message.SenderID), fmt.Sprintf("brrrrr!!! You said %s: %s", whereYouSaid, e.Message.Content))
				if err != nil {
					panic(err)
				}

				if resp.IsError() {
					panic(resp.Msg())
				}

				log.Printf("Sent message: %s - %s\n", e.Message.SenderFullName, e.Message.Content)
			}
		}
	}
}
