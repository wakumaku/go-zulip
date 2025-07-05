package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/realtime"
	"github.com/wakumaku/go-zulip/realtime/events"
)

func main() {
	zulipEmail := os.Getenv("ZULIP_EMAIL")
	zulipAPIKey := os.Getenv("ZULIP_API_KEY")
	zulipSite := os.Getenv("ZULIP_SITE")

	insecureClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	z, err := zulip.NewClient(zulip.Credentials(zulipSite, zulipEmail, zulipAPIKey),
		zulip.WithHTTPClient(&insecureClient),
		// zulip.WithPrintRequestData(),
		// zulip.WithPrintRawResponse(),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	realtimeSvc := realtime.NewService(z)

	const maxBackoffRetrySeconds = 30

	backoffRetrySeconds := 1

	for {
		if backoffRetrySeconds > 1 {
			log.Printf("Reconnecting in %d seconds ...\n", backoffRetrySeconds)
			time.Sleep(time.Duration(backoffRetrySeconds) * time.Second)
		}

		log.Printf("Connecting ...\n")

		backoffRetrySeconds *= 2
		if backoffRetrySeconds > maxBackoffRetrySeconds {
			backoffRetrySeconds = maxBackoffRetrySeconds
		}

		q, err := realtimeSvc.RegisterEvetQueue(ctx,
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
			),
			realtime.AllPublicStreams(true),
			// zulip.NarrowEvents(zulip.Narrower{}.Add(zulip.Is, "mentioned")),
		)
		if err != nil {
			log.Printf("error registering event queue: %s", err)
			continue
		}

		if q.IsError() {
			log.Printf("%s: %s", q.Msg(), q.Code())
			continue
		}

		log.Printf("QueueId: %s", q.QueueId)
		log.Println("Waiting for events...")

		lastEventID := q.LastEventId
		for {
			evs, err := realtimeSvc.GetEventsEventQueue(ctx, q.QueueId, realtime.LastEventID(lastEventID))
			if err != nil {
				log.Printf("error getting events from queue: %s", err)
				break
			}

			for _, ev := range evs.Events {
				if ev.EventType() == events.MessageType {
					message := ev.(*events.Message)
					if message.Message.DisplayRecipient.IsChannel {
						log.Printf("#%s [%s]: %s", message.Message.DisplayRecipient.Channel, message.Message.SenderFullName, message.Message.Content)
					} else {
						var users []string
						for _, user := range message.Message.DisplayRecipient.Users {
							users = append(users, user.FullName)
						}

						log.Printf("@%s: %s", strings.Join(users, ", @"), message.Message.Content)
					}
				} else if ev.EventType() == events.AlertWordsType {
					alertWords := ev.(*events.AlertWords)
					log.Printf("!AlertWords ID: %d, Words: %s", alertWords.ID, alertWords.AlertWords)
				} else if ev.EventType() == events.RealmUserType {
					realmUser := ev.(*events.RealmUser)
					log.Printf("@RealmUser ID: %d, Op: %s, FullName: %s", realmUser.ID, realmUser.Op, realmUser.Person.FullName)
				} else if ev.EventType() == events.PresenceType {
					presence := ev.(*events.Presence)
					log.Printf("*Presence Email: %s, Status: %s", presence.Email, presence.Presence.Website.Status)
				} else if ev.EventType() == events.RealmEmojiType {
					realmEmoji := ev.(*events.RealmEmoji)
					log.Printf(":RealmEmoji event ID: %d", realmEmoji.ID)
					// list emojis
					for id, emoji := range realmEmoji.RealmEmoji {
						log.Printf("  %s: %s, %s", id, emoji.Name, emoji.SourceURL)
					}
				} else if ev.EventType() == events.TypingType {
					typing := ev.(*events.Typing)
					log.Printf("...typing: : %d . %s -> %s", typing.ID, typing.Op, typing.Sender.Email)
				} else {
					log.Printf("#%d %s", ev.EventID(), ev.EventType())
				}

				lastEventID = ev.EventID()
			}
		}

		backoffRetrySeconds = 1
	}
}
