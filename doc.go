// Package zulip provides a Go client library for interacting with the Zulip API.
//
// The client supports various operations including:
//   - Sending and receiving messages
//   - Managing channels and subscriptions
//   - User management
//   - Real-time event streaming
//   - File uploads
//   - Organization management
//
// Usage:
//
//	// Create a client with credentials
//	credentials := zulip.Credentials("https://your-zulip-server.com", "your-email@example.com", "your-api-key")
//	client, err := zulip.NewClient(credentials)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Send a message
//	msgSvc := messages.NewService(client)
//	resp, err := msgSvc.SendMessage(ctx, recipient.ToChannel("general"), "Hello, World!")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// For more examples, see the examples directory.
package zulip
