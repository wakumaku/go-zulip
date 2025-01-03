package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/channels"
	"github.com/wakumaku/go-zulip/messages"
	"github.com/wakumaku/go-zulip/realtime"
	"github.com/wakumaku/go-zulip/users"
	"golang.org/x/sync/errgroup"
)

func main() {
	emailA := os.Getenv("ZULIP_EMAIL_A")
	apiKeyA := os.Getenv("ZULIP_API_KEY_A")
	emailB := os.Getenv("ZULIP_EMAIL_B")
	apiKeyB := os.Getenv("ZULIP_API_KEY_B")

	if emailA == "" || apiKeyA == "" || emailB == "" || apiKeyB == "" {
		log.Fatal("ZULIP_EMAIL_A, ZULIP_API_KEY_A, ZULIP_EMAIL_B and ZULIP_API_KEY_B are required")
	}

	site := "https://localhost"

	// Create an insecure client because of self-signed certificate
	insecureClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Create a client for user A
	clientA, err := zulip.NewClient(site, emailA, apiKeyA, zulip.WithHTTPClient(&insecureClient))
	if err != nil {
		log.Fatalf("failed to create zulip client for User A: %v", err)
	}

	// Create a client for user B
	clientB, err := zulip.NewClient(site, emailB, apiKeyB, zulip.WithHTTPClient(&insecureClient))
	if err != nil {
		log.Fatalf("failed to create zulip client for User B: %v", err)
	}

	// Create bots instantiating the services needed to interact with Zulip's API
	botA := Bot{
		UserSVC:     users.NewService(clientA),
		ChannelSVC:  channels.NewService(clientA),
		MessageSVC:  messages.NewService(clientA),
		RealtimeSVC: realtime.NewService(clientA),
	}

	botB := Bot{
		UserSVC:     users.NewService(clientB),
		ChannelSVC:  channels.NewService(clientB),
		MessageSVC:  messages.NewService(clientB),
		RealtimeSVC: realtime.NewService(clientB),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGHUP)
	defer cancel()

	// Channel and topic for the conversation
	channel := "talking_bots"
	topic := time.Now().Format("20060102") // use the current date as topic, eg.: 20250102

	log.Printf("Starting talking bots on channel %s with topic %s", channel, topic)

	// Run the bots concurrently
	errGrp := errgroup.Group{}
	errGrp.Go(func() error {
		return botA.Run(ctx, channel, topic)
	})

	// wait a bit so the conversation between the bots is established
	time.Sleep(3 * time.Second)

	errGrp.Go(func() error {
		return botB.Run(ctx, channel, topic)
	})

	if err := errGrp.Wait(); err != nil {
		log.Fatalf("failed to run talking bots: %v", err)
	}

	log.Println("Shutting down talking bots")
}
