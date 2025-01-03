package integration

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/channels"
	"github.com/wakumaku/go-zulip/invitations"
	"github.com/wakumaku/go-zulip/messages"
	"github.com/wakumaku/go-zulip/org"
	"github.com/wakumaku/go-zulip/realtime"
	"github.com/wakumaku/go-zulip/realtime/events"
	"github.com/wakumaku/go-zulip/specialty"
	"github.com/wakumaku/go-zulip/users"
)

func TestIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	zulipEmail := os.Getenv("ZULIP_EMAIL")
	zulipAPIKey := os.Getenv("ZULIP_API_KEY")
	zulipSite := os.Getenv("ZULIP_SITE")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	insecureClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	adminClient, err := zulip.NewClient(zulipSite, zulipEmail, zulipAPIKey,
		zulip.WithHTTPClient(&insecureClient),
		// zulip.WithPrintRequestData(),
		// zulip.WithPrintRawResponse(),
	)
	assert.NoError(t, err)

	adminUserSvc := users.NewService(adminClient)
	adminSpecialtySvc := specialty.NewService(adminClient)

	userPassword := "1234567890ABcd"
	userSuffix := time.Now().UTC().Format("150405")
	// create user A
	userAEmail := fmt.Sprintf("usera_%s@zulip.test", userSuffix)
	respCreateUserA, err := adminUserSvc.CreateUser(ctx, userAEmail, userPassword, "User A")
	assert.NoError(t, err)
	assert.True(t, respCreateUserA.IsSuccess())

	userAID := respCreateUserA.UserID

	// get API Key User A
	respFetchAPIKeyA, err := adminSpecialtySvc.FetchAPIKeyProduction(ctx, userAEmail, userPassword)
	assert.NoError(t, err)
	assert.True(t, respFetchAPIKeyA.IsSuccess())
	// User A Client
	userA, err := zulip.NewClient(zulipSite, userAEmail, respFetchAPIKeyA.APIKey,
		zulip.WithHTTPClient(&insecureClient),
	)
	assert.NoError(t, err)

	// User A sends a message
	userAMsgSvc := messages.NewService(userA)
	userASendMsgResp, err := userAMsgSvc.SendMessage(ctx, messages.ToChannelTopic(messages.ToChannelName("general"), "greetings"), "Im User A!")
	assert.NoError(t, err)
	assert.True(t, userASendMsgResp.IsSuccess())

	// create user B
	userBEmail := fmt.Sprintf("userb_%s@zulip.test", userSuffix)
	respCreateUserB, err := adminUserSvc.CreateUser(ctx, userBEmail, userPassword, "User B")
	assert.NoError(t, err)
	assert.True(t, respCreateUserB.IsSuccess())

	// get API Key User B
	respFetchAPIKeyB, err := adminSpecialtySvc.FetchAPIKeyProduction(ctx, userBEmail, userPassword)
	assert.NoError(t, err)
	assert.True(t, respFetchAPIKeyB.IsSuccess())
	// User B Client
	userB, err := zulip.NewClient(zulipSite, userBEmail, respFetchAPIKeyB.APIKey,
		zulip.WithHTTPClient(&insecureClient),
		// zulip.WithPrintRequestData(),
		// zulip.WithPrintRawResponse(),
	)
	assert.NoError(t, err)
	adminInvitationSvc := invitations.NewService(adminClient)
	respCreateReusableLink, err := adminInvitationSvc.CreateReusableInvitationLink(ctx,
		invitations.IncludeRealmDefaultSubscriptions(true),
		invitations.InviteAs(zulip.RoleMember),
		invitations.InviteExpiresInMinutes(15),
		invitations.StreamIds([]int{1, 2}),
	)
	assert.NoError(t, err)
	assert.True(t, respCreateReusableLink.IsSuccess())

	adminMsgSvc := messages.NewService(adminClient)
	respSendMessage, err := adminMsgSvc.SendMessage(ctx, messages.ToChannelTopic(messages.ToChannelName("general"), "greetings"), "jello güorld")
	assert.NoError(t, err)
	assert.Equal(t, respSendMessage.HTTPCode(), http.StatusOK)
	assert.Equal(t, respSendMessage.Result(), zulip.ResultSuccess)

	respSendMessage, err = adminMsgSvc.SendMessage(ctx, messages.ToChannelName("nonexistent"), "jello güorld", messages.ToTopic("greetings"))
	assert.NoError(t, err)
	assert.Equal(t, respSendMessage.HTTPCode(), http.StatusBadRequest)
	assert.Equal(t, respSendMessage.Result(), zulip.ResultError)
	assert.Equal(t, respSendMessage.Code(), "STREAM_DOES_NOT_EXIST")

	respSendMessage, err = adminMsgSvc.SendMessage(ctx, messages.ToUserName("eeshan@zulip.com"), "jello güorld")
	assert.NoError(t, err)
	assert.Equal(t, respSendMessage.HTTPCode(), http.StatusBadRequest)
	assert.Equal(t, respSendMessage.Result(), zulip.ResultError)
	assert.Equal(t, respSendMessage.Code(), "BAD_REQUEST")
	assert.Equal(t, respSendMessage.Msg(), "Invalid email 'eeshan@zulip.com'")

	// Send message with a picture
	fileToUpload := "zulip-desktop-screenshot.webp"
	respUploadFile, err := adminMsgSvc.UploadFile(ctx, "../testdata/"+fileToUpload)
	assert.NoError(t, err)
	assert.Equal(t, respUploadFile.HTTPCode(), http.StatusOK)
	assert.Equal(t, respUploadFile.Result(), zulip.ResultSuccess)
	uploadedFileURL, err := respUploadFile.FieldValue("url")
	assert.NoError(t, err)
	assert.Contains(t, uploadedFileURL, fileToUpload)

	messageWithPicture := fmt.Sprintf("Here a picture: [this picture](%s) of my castle!", uploadedFileURL)
	respSendMessageWithPicture, err := adminMsgSvc.SendMessage(ctx, messages.ToChannelTopic(messages.ToChannelName("general"), "greetings"), messageWithPicture)
	assert.NoError(t, err)
	assert.Equal(t, respSendMessageWithPicture.HTTPCode(), http.StatusOK)
	assert.Equal(t, respSendMessageWithPicture.Result(), zulip.ResultSuccess)
	assert.NoError(t, err)

	fileToUpload2 := "zulip-mobile-screenshot.jpg"
	respUploadFile2, err := adminMsgSvc.UploadFile(ctx, "../testdata/"+fileToUpload2)
	assert.NoError(t, err)
	assert.Equal(t, respUploadFile2.HTTPCode(), http.StatusOK)
	assert.Equal(t, respUploadFile2.Result(), zulip.ResultSuccess)
	uploadedFileURL2, err := respUploadFile2.FieldValue("url")
	assert.NoError(t, err)
	assert.Contains(t, uploadedFileURL2, fileToUpload2)

	messageWithPicture2 := fmt.Sprintf("This is an image [%s](%s) :camera:", fileToUpload2, uploadedFileURL2)
	respSendMessageWithPicture2, err := adminMsgSvc.SendMessage(ctx, messages.ToChannelTopic(messages.ToChannelName("general"), "greetings"), messageWithPicture2)
	assert.NoError(t, err)
	assert.Equal(t, respSendMessageWithPicture2.HTTPCode(), http.StatusOK)
	assert.Equal(t, respSendMessageWithPicture2.Result(), zulip.ResultSuccess)
	assert.NoError(t, err)

	// React to the message
	addEmojiReaction, err := adminMsgSvc.AddEmojiReaction(ctx, respSendMessageWithPicture2.ID, "tada")
	assert.NoError(t, err)
	assert.True(t, addEmojiReaction.IsSuccess())

	addEmojiReaction, err = adminMsgSvc.AddEmojiReaction(ctx, respSendMessageWithPicture2.ID, "heart")
	assert.NoError(t, err)
	assert.True(t, addEmojiReaction.IsSuccess())

	// Remove reaction to a message
	// 1. Send a reaction
	addEmojiReactionToRemove, err := adminMsgSvc.AddEmojiReaction(ctx, respSendMessageWithPicture.ID, "+1")
	assert.NoError(t, err)
	assert.True(t, addEmojiReactionToRemove.IsSuccess())

	removeEmojiReaction, err := adminMsgSvc.RemoveEmojiReaction(ctx, respSendMessageWithPicture.ID, messages.RemoveEmojiReactionEmojiName("+1"))
	assert.NoError(t, err)
	assert.True(t, removeEmojiReaction.IsSuccess())

	// Edit the message with the picture
	respEditMessage, err := adminMsgSvc.EditMessage(ctx, respSendMessageWithPicture.ID,
		messages.NewContent("Message EDITED: :pencil::"+messageWithPicture),
		messages.SetPropagateMode(messages.PropagateModeAll),
		messages.SendNotificationToNewThread(true),
		messages.SendNotificationToOldThread(true),
		messages.MoveToTopic("pictures"),
	)
	assert.NoError(t, err)
	assert.Equal(t, respEditMessage.HTTPCode(), http.StatusOK)
	assert.Equal(t, respEditMessage.Result(), zulip.ResultSuccess)

	// Reedit the message, but no changes
	respEditMessage, err = adminMsgSvc.EditMessage(ctx, respSendMessageWithPicture.ID)
	assert.NoError(t, err)
	assert.Equal(t, respEditMessage.HTTPCode(), http.StatusBadRequest)
	assert.Equal(t, respEditMessage.Result(), zulip.ResultError)
	assert.Equal(t, respEditMessage.Msg(), "Nothing to change")

	// Delete the message
	// 1. Send the message
	respSendMessageToDelete, err := adminMsgSvc.SendMessage(ctx, messages.ToChannelTopic(messages.ToChannelName("general"), "greetings"), "this message will be deleted")
	assert.NoError(t, err)
	assert.Equal(t, respSendMessageToDelete.HTTPCode(), http.StatusOK)
	assert.Equal(t, respSendMessageToDelete.Result(), zulip.ResultSuccess)
	// 2. Delete
	respDeleteMessage, err := adminMsgSvc.DeleteMessage(ctx, respSendMessageToDelete.ID)
	assert.NoError(t, err)
	assert.Equal(t, respDeleteMessage.HTTPCode(), http.StatusOK)
	assert.Equal(t, respDeleteMessage.Result(), zulip.ResultSuccess)

	// registering for events
	adminRealtimeSvc := realtime.NewService(adminClient)
	respRegister, err := adminRealtimeSvc.RegisterEvetQueue(ctx, realtime.EventTypes(events.MessageType))
	assert.NoError(t, err)
	assert.True(t, respRegister.IsSuccess())

	// USER A sends some messages...
	for range 1 {
		respSendMessageEvent, err := userAMsgSvc.SendMessage(ctx,
			messages.ToChannelTopic(messages.ToChannelName("general"), "greetings"),
			"Im USER A Sending a message so its captured by a registered queue",
		)

		assert.NoError(t, err)
		assert.Equal(t, respSendMessageEvent.HTTPCode(), http.StatusOK)
		assert.Equal(t, respSendMessageEvent.Result(), zulip.ResultSuccess)
		time.Sleep(100 * time.Millisecond)
	}

	for range 1 {
		respSendMessageEvent, err := adminMsgSvc.SendMessage(ctx,
			messages.ToChannelTopic(messages.ToChannelName("general"), "greetings"),
			"Im ADMIN Sending a message so its captured by a registered queue",
		)

		assert.NoError(t, err)
		assert.Equal(t, respSendMessageEvent.HTTPCode(), http.StatusOK)
		assert.Equal(t, respSendMessageEvent.Result(), zulip.ResultSuccess)
		time.Sleep(100 * time.Millisecond)
	}

	// receive the message via event
	respGetEvents, err := adminRealtimeSvc.GetEventsEventQueue(ctx, respRegister.QueueId)
	assert.NoError(t, err)
	assert.Equal(t, respGetEvents.HTTPCode(), http.StatusOK)
	assert.Equal(t, respGetEvents.Result(), zulip.ResultSuccess)

	// Send a message and get it later
	messageToBeGetLaterMessage := "Message to be get later"
	messageToBeGetLater, err := userAMsgSvc.SendMessage(ctx,
		messages.ToChannelTopic(messages.ToChannelName("general"), "greetings"),
		messageToBeGetLaterMessage,
	)

	assert.NoError(t, err)
	assert.Equal(t, messageToBeGetLater.HTTPCode(), http.StatusOK)
	assert.Equal(t, messageToBeGetLater.Result(), zulip.ResultSuccess)

	// Get the message
	respGetMessage, err := userAMsgSvc.GetMessages(ctx,
		messages.Anchor("newest"),
		messages.NumBefore(1),
		messages.NumAfter(1),
		messages.NarrowMessage(zulip.Narrower{}.
			Add(zulip.Id, messageToBeGetLater.ID),
		),
		messages.ApplyMarkdownMessage(false),
	)
	assert.NoError(t, err)
	assert.Equal(t, respGetMessage.HTTPCode(), http.StatusOK)

	assert.Equal(t, messageToBeGetLater.ID, respGetMessage.Messages[0].ID)
	assert.Equal(t, messageToBeGetLaterMessage, respGetMessage.Messages[0].Content)
	assert.Equal(t, "greetings", respGetMessage.Messages[0].Subject)
	assert.Equal(t, "general", respGetMessage.Messages[0].DisplayRecipient.Channel)

	// Testing narrower with multiple conditions
	narrowMultiple := zulip.Narrower{}.
		Add(zulip.Channel, "general").
		AddNegated(zulip.IsUnread, nil).
		Add(zulip.Search, `"Message to be get later"`)

	respGetMessage, err = userAMsgSvc.GetMessages(ctx, messages.Anchor("newest"),
		messages.NumBefore(1),
		messages.NumAfter(1),
		messages.NarrowMessage(narrowMultiple),
		messages.ApplyMarkdownMessage(false),
	)

	assert.NoError(t, err)
	assert.Equal(t, respGetMessage.HTTPCode(), http.StatusOK)
	assert.Equal(t, respGetMessage.Result(), zulip.ResultSuccess)

	// fetch a single message
	respFetchSingleMessage, err := userAMsgSvc.FetchSingleMessage(ctx, respGetMessage.Messages[0].ID, messages.ApplyMarkdownSingleMessage(false))
	assert.NoError(t, err)
	assert.Equal(t, respFetchSingleMessage.HTTPCode(), http.StatusOK)
	assert.Equal(t, respFetchSingleMessage.Result(), zulip.ResultSuccess)

	assert.Equal(t, respFetchSingleMessage.Message.ID, respGetMessage.Messages[0].ID)
	assert.Equal(t, respFetchSingleMessage.Message.Content, respGetMessage.Messages[0].Content)

	// Channel subscriptions
	adminChannels := channels.NewService(adminClient)
	respGetSubscribedChannels, err := adminChannels.GetSubscribedChannels(ctx)
	assert.NoError(t, err)
	assert.Equal(t, respGetSubscribedChannels.HTTPCode(), http.StatusOK)
	assert.Equal(t, respGetSubscribedChannels.Result(), zulip.ResultSuccess)

	// the user is subscribed to the general channel by default
	generalChannelFound := false
	for _, subscription := range respGetSubscribedChannels.Subscriptions {
		if subscription.Name == "general" {
			generalChannelFound = true
			break
		}
	}
	assert.True(t, generalChannelFound, "User should be subscribed to the general channel by default")

	// subscribe to a new channel and resubscribe to the general channel (just to test the API response containing more info)
	newChannelName := "new_channel_" + time.Now().UTC().Format("150405")
	respSubscribeToChannel, err := adminChannels.SubscribeToChannel(ctx, []channels.SubscribeTo{{Name: newChannelName, Description: "da cool name"}, {Name: "general"}})
	assert.NoError(t, err)
	assert.Equal(t, respSubscribeToChannel.HTTPCode(), http.StatusOK)
	assert.Equal(t, respSubscribeToChannel.Result(), zulip.ResultSuccess)

	userAChannels := channels.NewService(userA)
	respSubscribeToChannelUserA, err := userAChannels.SubscribeToChannel(ctx, []channels.SubscribeTo{{Name: newChannelName, Description: "da cool name"}, {Name: "general"}})
	assert.NoError(t, err)
	assert.Equal(t, respSubscribeToChannelUserA.HTTPCode(), http.StatusOK)
	assert.Equal(t, respSubscribeToChannelUserA.Result(), zulip.ResultSuccess)

	respSendMessageNewChannel, err := userAMsgSvc.SendMessage(ctx, messages.ToChannelName(newChannelName), "Hello new channel!", messages.ToTopic("ThatsNew"))
	assert.NoError(t, err)
	assert.Equal(t, respSendMessageNewChannel.HTTPCode(), http.StatusOK)
	assert.Equal(t, respSendMessageNewChannel.Result(), zulip.ResultSuccess)

	// validate the subscription
	respGetSubscribedChannels, err = adminChannels.GetSubscribedChannels(ctx, channels.IncludeSubscribersList(true))
	assert.NoError(t, err)
	assert.Equal(t, respGetSubscribedChannels.HTTPCode(), http.StatusOK)
	assert.Equal(t, respGetSubscribedChannels.Result(), zulip.ResultSuccess)

	// the user is subscribed to the newChannelName too
	newChannelFound := false
	for _, subscription := range respGetSubscribedChannels.Subscriptions {
		if subscription.Name == newChannelName {
			newChannelFound = true
			break
		}
	}
	assert.True(t, newChannelFound, "User should be subscribed to the "+newChannelName+" channel")

	// 3 People private chat: Admin, User A, User B
	// Admin sends a message to User A and User B
	respCreatePrivateChat, err := adminMsgSvc.SendMessage(ctx, messages.ToUserNames([]string{userAEmail, userBEmail}), "Hello User A and User B!")
	assert.NoError(t, err)
	assert.Equal(t, respCreatePrivateChat.HTTPCode(), http.StatusOK)
	assert.Equal(t, respCreatePrivateChat.Result(), zulip.ResultSuccess)
	// User A sends a message to Admin and User B
	respCreatePrivateChat, err = userAMsgSvc.SendMessage(ctx, messages.ToUserNames([]string{zulipEmail, userBEmail}), "Hello Admin and User B!")
	assert.NoError(t, err)
	assert.Equal(t, respCreatePrivateChat.HTTPCode(), http.StatusOK)
	assert.Equal(t, respCreatePrivateChat.Result(), zulip.ResultSuccess)
	// User B sends a message to Admin and User A
	userBMsgSvc := messages.NewService(userB)
	respCreatePrivateChat, err = userBMsgSvc.SendMessage(ctx, messages.ToUserNames([]string{zulipEmail, userAEmail}), "Hello Admin and User A!")
	assert.NoError(t, err)
	assert.Equal(t, respCreatePrivateChat.HTTPCode(), http.StatusOK)
	assert.Equal(t, respCreatePrivateChat.Result(), zulip.ResultSuccess)

	// UserA gets its own information
	uaserAUserSvc := users.NewService(userA)
	respGetUserMe, err := uaserAUserSvc.GetUserMe(ctx)
	assert.NoError(t, err)
	assert.Equal(t, respGetUserMe.HTTPCode(), http.StatusOK)
	assert.Equal(t, respGetUserMe.Result(), zulip.ResultSuccess)
	assert.Equal(t, respGetUserMe.Email, userAEmail)

	// Get User A information by ID
	respGetUser, err := adminUserSvc.GetUser(ctx, respGetUserMe.UserID)
	assert.NoError(t, err)
	assert.Equal(t, respGetUser.HTTPCode(), http.StatusOK)
	assert.Equal(t, respGetUser.Result(), zulip.ResultSuccess)
	assert.Equal(t, userAEmail, respGetUser.User.Email)

	// Get ALL Users
	respGetUsers, err := adminUserSvc.GetUsers(ctx,
		users.ClientGravatars(true),
		users.IncludeCustomProfilesFields(true),
	)
	assert.NoError(t, err)
	assert.Equal(t, respGetUsers.HTTPCode(), http.StatusOK)
	assert.Equal(t, respGetUsers.Result(), zulip.ResultSuccess)
	assert.GreaterOrEqual(t, len(respGetUsers.Members), 3)

	// Status Text and Emoji
	// Set status text and emoji
	userAStatusText := "I'm busy"
	respUpdateStatus, err := uaserAUserSvc.UpdateStatus(ctx,
		users.StatusText(userAStatusText),
		users.StatusEmojiName("thumbs_up"),
		// zulip.StatusEmojiCode("1f389"),
		// zulip.StatusReactionType(zulip.ReactionTypeUnicode),
	)
	assert.NoError(t, err)
	assert.Equal(t, respUpdateStatus.HTTPCode(), http.StatusOK)
	assert.Equal(t, respUpdateStatus.Result(), zulip.ResultSuccess)

	// Get User A status
	respGetUserStatus, err := adminUserSvc.GetUserStatus(ctx, userAID)
	assert.NoError(t, err)
	assert.Equal(t, respGetUserStatus.HTTPCode(), http.StatusOK)
	assert.Equal(t, respGetUserStatus.Result(), zulip.ResultSuccess)
	assert.Equal(t, respGetUserStatus.Status.StatusText, userAStatusText)

	// Upload a custom emoji and set it as the user's status emoji
	adminOrgSvc := org.NewService(adminClient)
	emojiName := "dancing_gopher_" + userSuffix
	respUploadEmoji, err := adminOrgSvc.UploadCustomEmoji(ctx, emojiName, "../testdata/dancing_gopher.gif")
	assert.NoError(t, err)
	assert.Equal(t, respUploadEmoji.HTTPCode(), http.StatusOK)
	assert.Equal(t, respUploadEmoji.Result(), zulip.ResultSuccess)

	// Set the custom emoji as the user's status emoji
	respUpdateStatus, err = uaserAUserSvc.UpdateStatus(ctx,
		users.StatusText("I'm dancing"),
		users.StatusEmojiName(emojiName),
	)
	assert.NoError(t, err)
	assert.Equal(t, respUpdateStatus.HTTPCode(), http.StatusOK)
	assert.Equal(t, respUpdateStatus.Result(), zulip.ResultSuccess)

	// spew.Dump("ok")
}
