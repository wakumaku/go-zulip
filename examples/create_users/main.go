package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/specialty"
	"github.com/wakumaku/go-zulip/users"
)

func main() {
	// Must be an admin with capabilities to create users
	email := os.Getenv("ZULIP_EMAIL")
	apiKey := os.Getenv("ZULIP_API_KEY")
	site := os.Getenv("ZULIP_SITE")

	var (
		userEmail    string
		userPassword string
		userName     string
	)
	flag.StringVar(&userEmail, "email", "", "email of the user to create")
	flag.StringVar(&userPassword, "password", "", "password of the user to create")
	flag.StringVar(&userName, "name", "", "name of the user to create")
	flag.Parse()

	if userEmail == "" || userPassword == "" || userName == "" {
		log.Fatal("email, password and name are required")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	insecureClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	admin, err := zulip.NewClient(zulip.Credentials(site, email, apiKey),
		zulip.WithHTTPClient(&insecureClient),
	)
	if err != nil {
		log.Fatalf("failed to create zulip client: %v", err)
	}

	// Service to create users
	usersSvc := users.NewService(admin)
	// Create the user
	createUserResp, err := usersSvc.CreateUser(ctx, userEmail, userPassword, userName)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	if createUserResp.IsError() {
		log.Fatalf("zulip API error creating user: %v", createUserResp.Msg())
	}

	// Service to fetch API Key
	specialtySvc := specialty.NewService(admin)
	// Fetch the API Key
	fetchAPIKeyResp, err := specialtySvc.FetchAPIKeyProduction(ctx, userEmail, userPassword)
	if err != nil {
		log.Fatalf("failed to fetch API Key: %v", err)
	}

	if fetchAPIKeyResp.IsError() {
		log.Fatalf("zulip API error fetching API Key: %v", fetchAPIKeyResp.Msg())
	}

	log.Println("User created successfully")
	fmt.Printf("\tEmail: %s\n", userEmail)
	fmt.Printf("\tPassword: %s\n", userPassword)
	fmt.Printf("\tAPI Key: %s\n", fetchAPIKeyResp.APIKey)
	fmt.Println("Done.")
}
