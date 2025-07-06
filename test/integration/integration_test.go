package integration

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/channels"
	"github.com/wakumaku/go-zulip/invitations"
	"github.com/wakumaku/go-zulip/messages"
	"github.com/wakumaku/go-zulip/messages/recipient"
	"github.com/wakumaku/go-zulip/narrow"
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

	if zulipEmail == "" || zulipAPIKey == "" || zulipSite == "" {
		t.Skip("ZULIP_EMAIL, ZULIP_API_KEY, and ZULIP_SITE environment variables must be set for integration tests")
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

	debugLogger := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	adminClient, err := zulip.NewClient(zulip.Credentials(zulipSite, zulipEmail, zulipAPIKey),
		zulip.WithHTTPClient(&insecureClient),
		zulip.WithLogger(debugLogger),
	)
	require.NoError(t, err)

	adminUserSvc := users.NewService(adminClient)
	adminSpecialtySvc := specialty.NewService(adminClient)

	userPassword := "1234567890ABcd"
	userSuffix := time.Now().UTC().Format("150405")
	// create user A
	userAEmail := fmt.Sprintf("usera_%s@zulip.test", userSuffix)
	respCreateUserA, err := adminUserSvc.CreateUser(ctx, userAEmail, userPassword, "User A")
	require.NoError(t, err)
	assert.True(t, respCreateUserA.IsSuccess())

	userAID := respCreateUserA.UserID

	// get API Key User A
	respFetchAPIKeyA, err := adminSpecialtySvc.FetchAPIKeyProduction(ctx, userAEmail, userPassword)
	require.NoError(t, err)
	assert.True(t, respFetchAPIKeyA.IsSuccess())
	// User A Client
	userA, err := zulip.NewClient(zulip.Credentials(zulipSite, userAEmail, respFetchAPIKeyA.APIKey),
		zulip.WithHTTPClient(&insecureClient),
		zulip.WithLogger(debugLogger),
	)
	require.NoError(t, err)

	// User A sends a message
	userAMsgSvc := messages.NewService(userA)
	userASendMsgResp, err := userAMsgSvc.SendMessageToChannelTopic(ctx, recipient.ToChannel("general"), "greetings", "Im User A!")
	require.NoError(t, err)
	assert.True(t, userASendMsgResp.IsSuccess())

	// create user B
	userBEmail := fmt.Sprintf("userb_%s@zulip.test", userSuffix)
	respCreateUserB, err := adminUserSvc.CreateUser(ctx, userBEmail, userPassword, "User B")
	require.NoError(t, err)
	assert.True(t, respCreateUserB.IsSuccess())

	// get API Key User B
	respFetchAPIKeyB, err := adminSpecialtySvc.FetchAPIKeyProduction(ctx, userBEmail, userPassword)
	require.NoError(t, err)
	assert.True(t, respFetchAPIKeyB.IsSuccess())
	// User B Client
	userB, err := zulip.NewClient(zulip.Credentials(zulipSite, userBEmail, respFetchAPIKeyB.APIKey),
		zulip.WithHTTPClient(&insecureClient),
		zulip.WithLogger(debugLogger),
	)
	require.NoError(t, err)

	adminInvitationSvc := invitations.NewService(adminClient)
	respCreateReusableLink, err := adminInvitationSvc.CreateReusableInvitationLink(ctx,
		invitations.IncludeRealmDefaultSubscriptions(true),
		invitations.InviteAs(zulip.MemberRole),
		invitations.InviteExpiresInMinutes(15),
		invitations.StreamIds([]int{1, 2}),
	)
	require.NoError(t, err)
	assert.True(t, respCreateReusableLink.IsSuccess())

	adminMsgSvc := messages.NewService(adminClient)
	respSendMessage, err := adminMsgSvc.SendMessageToChannelTopic(ctx, recipient.ToChannel("general"), "greetings", "jello güorld")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respSendMessage.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respSendMessage.Result())

	respSendMessage, err = adminMsgSvc.SendMessage(ctx, recipient.ToChannel("nonexistent"), "jello güorld", messages.ToTopic("greetings"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, respSendMessage.HTTPCode())
	assert.Equal(t, zulip.ResultError, respSendMessage.Result())
	assert.Equal(t, "STREAM_DOES_NOT_EXIST", respSendMessage.Code())

	respSendMessage, err = adminMsgSvc.SendMessage(ctx, recipient.ToUser("eeshan@zulip.com"), "jello güorld")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, respSendMessage.HTTPCode())
	assert.Equal(t, zulip.ResultError, respSendMessage.Result())
	assert.Equal(t, "BAD_REQUEST", respSendMessage.Code())
	assert.Equal(t, "Invalid email 'eeshan@zulip.com'", respSendMessage.Msg())

	// Send message with a picture
	fileToUpload := "zulip-desktop-screenshot.webp"
	respUploadFile, err := adminMsgSvc.UploadFile(ctx, "../testdata/"+fileToUpload)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUploadFile.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respUploadFile.Result())
	uploadedFileURL, err := respUploadFile.FieldValue("url")
	require.NoError(t, err)
	assert.Contains(t, uploadedFileURL, fileToUpload)

	messageWithPicture := fmt.Sprintf("Here a picture: [this picture](%s) of my castle!", uploadedFileURL)
	respSendMessageWithPicture, err := adminMsgSvc.SendMessageToChannelTopic(ctx, recipient.ToChannel("general"), "greetings", messageWithPicture)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respSendMessageWithPicture.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respSendMessageWithPicture.Result())
	require.NoError(t, err)

	fileToUpload2 := "zulip-mobile-screenshot.jpg"
	respUploadFile2, err := adminMsgSvc.UploadFile(ctx, "../testdata/"+fileToUpload2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUploadFile2.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respUploadFile2.Result())
	uploadedFileURL2, err := respUploadFile2.FieldValue("url")
	require.NoError(t, err)
	assert.Contains(t, uploadedFileURL2, fileToUpload2)

	messageWithPicture2 := fmt.Sprintf("This is an image [%s](%s) :camera:", fileToUpload2, uploadedFileURL2)
	respSendMessageWithPicture2, err := adminMsgSvc.SendMessageToChannelTopic(ctx, recipient.ToChannel("general"), "greetings", messageWithPicture2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respSendMessageWithPicture2.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respSendMessageWithPicture2.Result())
	require.NoError(t, err)

	// React to the message
	addEmojiReaction, err := adminMsgSvc.AddEmojiReaction(ctx, respSendMessageWithPicture2.ID, "tada")
	require.NoError(t, err)
	assert.True(t, addEmojiReaction.IsSuccess())

	addEmojiReaction, err = adminMsgSvc.AddEmojiReaction(ctx, respSendMessageWithPicture2.ID, "heart")
	require.NoError(t, err)
	assert.True(t, addEmojiReaction.IsSuccess())

	// Remove reaction to a message
	// 1. Send a reaction
	addEmojiReactionToRemove, err := adminMsgSvc.AddEmojiReaction(ctx, respSendMessageWithPicture.ID, "+1")
	require.NoError(t, err)
	assert.True(t, addEmojiReactionToRemove.IsSuccess())

	removeEmojiReaction, err := adminMsgSvc.RemoveEmojiReaction(ctx, respSendMessageWithPicture.ID, messages.RemoveEmojiReactionEmojiName("+1"))
	require.NoError(t, err)
	assert.True(t, removeEmojiReaction.IsSuccess())

	// Edit the message with the picture
	respEditMessage, err := adminMsgSvc.EditMessage(ctx, respSendMessageWithPicture.ID,
		messages.NewContent("Message EDITED: :pencil::"+messageWithPicture),
		messages.SetPropagateMode(messages.PropagateModeAll),
		messages.SendNotificationToNewThread(true),
		messages.SendNotificationToOldThread(true),
		messages.MoveToTopic("pictures"),
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respEditMessage.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respEditMessage.Result())

	// Reedit the message, but no changes
	respEditMessage, err = adminMsgSvc.EditMessage(ctx, respSendMessageWithPicture.ID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, respEditMessage.HTTPCode())
	assert.Equal(t, zulip.ResultError, respEditMessage.Result())
	assert.Equal(t, "Nothing to change", respEditMessage.Msg())

	// Delete the message
	// 1. Send the message
	respSendMessageToDelete, err := adminMsgSvc.SendMessageToChannelTopic(ctx, recipient.ToChannel("general"), "greetings", "this message will be deleted")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respSendMessageToDelete.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respSendMessageToDelete.Result())
	// 2. Delete
	respDeleteMessage, err := adminMsgSvc.DeleteMessage(ctx, respSendMessageToDelete.ID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respDeleteMessage.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respDeleteMessage.Result())

	// registering for events
	adminRealtimeSvc := realtime.NewService(adminClient)
	respRegister, err := adminRealtimeSvc.RegisterEvetQueue(ctx, realtime.EventTypes(events.MessageType))
	require.NoError(t, err)
	assert.True(t, respRegister.IsSuccess())

	// USER A sends some messages...
	for range 1 {
		respSendMessageEvent, err := userAMsgSvc.SendMessageToChannelTopic(ctx,
			recipient.ToChannel("general"), "greetings",
			"Im USER A Sending a message so its captured by a registered queue",
		)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, respSendMessageEvent.HTTPCode())
		assert.Equal(t, zulip.ResultSuccess, respSendMessageEvent.Result())
		time.Sleep(100 * time.Millisecond)
	}

	for range 1 {
		respSendMessageEvent, err := adminMsgSvc.SendMessageToChannelTopic(ctx,
			recipient.ToChannel("general"), "greetings",
			"Im ADMIN Sending a message so its captured by a registered queue",
		)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, respSendMessageEvent.HTTPCode())
		assert.Equal(t, zulip.ResultSuccess, respSendMessageEvent.Result())
		time.Sleep(100 * time.Millisecond)
	}

	// receive the message via event
	respGetEvents, err := adminRealtimeSvc.GetEventsEventQueue(ctx, respRegister.QueueID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetEvents.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetEvents.Result())

	// Send a message and get it later
	messageToBeGetLaterMessage := "Message to be get later"
	messageToBeGetLater, err := userAMsgSvc.SendMessageToChannelTopic(ctx,
		recipient.ToChannel("general"), "greetings",
		messageToBeGetLaterMessage,
	)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, messageToBeGetLater.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, messageToBeGetLater.Result())

	// Get the message
	respGetMessage, err := userAMsgSvc.GetMessages(ctx,
		messages.Anchor("newest"),
		messages.NumBefore(1),
		messages.NumAfter(1),
		messages.NarrowMessage(narrow.NewFilter().
			Add(narrow.New(narrow.ID, messageToBeGetLater.ID)),
		),
		messages.ApplyMarkdownMessage(false),
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetMessage.HTTPCode())

	assert.Equal(t, messageToBeGetLater.ID, respGetMessage.Messages[0].ID)
	assert.Equal(t, messageToBeGetLaterMessage, respGetMessage.Messages[0].Content)
	assert.Equal(t, "greetings", respGetMessage.Messages[0].Subject)
	assert.Equal(t, "general", respGetMessage.Messages[0].DisplayRecipient.Channel)

	// Testing narrower with multiple conditions
	narrowMultiple := narrow.NewFilter().
		Add(narrow.New(narrow.Channel, "general")).
		Add(narrow.New(narrow.Topic, "greetings")).
		// Add(narrow.Negate(narrow.IsUnread)).
		Add(narrow.New(narrow.Search, `Message to be get later`))

	respGetMessageNarrow, err := userAMsgSvc.GetMessages(ctx, messages.Anchor("newest"),
		messages.NumBefore(1),
		messages.NumAfter(1),
		messages.NarrowMessage(narrowMultiple),
		messages.ApplyMarkdownMessage(false),
	)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetMessageNarrow.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetMessageNarrow.Result())

	assert.Equal(t, messageToBeGetLater.ID, respGetMessageNarrow.Messages[0].ID)
	assert.Equal(t, messageToBeGetLaterMessage, respGetMessageNarrow.Messages[0].Content)
	assert.Equal(t, "general", respGetMessageNarrow.Messages[0].DisplayRecipient.Channel)
	assert.Equal(t, "greetings", respGetMessageNarrow.Messages[0].Subject)

	// fetch a single message
	respFetchSingleMessage, err := userAMsgSvc.FetchSingleMessage(ctx, respGetMessage.Messages[0].ID, messages.ApplyMarkdownSingleMessage(false))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respFetchSingleMessage.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respFetchSingleMessage.Result())

	assert.Equal(t, respFetchSingleMessage.Message.ID, respGetMessage.Messages[0].ID)
	assert.Equal(t, respFetchSingleMessage.Message.Content, respGetMessage.Messages[0].Content)

	// Channel subscriptions
	adminChannels := channels.NewService(adminClient)
	respGetSubscribedChannels, err := adminChannels.GetSubscribedChannels(ctx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetSubscribedChannels.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetSubscribedChannels.Result())

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
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respSubscribeToChannel.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respSubscribeToChannel.Result())

	userAChannels := channels.NewService(userA)
	respSubscribeToChannelUserA, err := userAChannels.SubscribeToChannel(ctx, []channels.SubscribeTo{{Name: newChannelName, Description: "da cool name"}, {Name: "general"}})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respSubscribeToChannelUserA.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respSubscribeToChannelUserA.Result())

	respSendMessageNewChannel, err := userAMsgSvc.SendMessage(ctx, recipient.ToChannel(newChannelName), "Hello new channel!", messages.ToTopic("ThatsNew"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respSendMessageNewChannel.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respSendMessageNewChannel.Result())

	// validate the subscription
	respGetSubscribedChannels, err = adminChannels.GetSubscribedChannels(ctx, channels.IncludeSubscribersList(true))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetSubscribedChannels.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetSubscribedChannels.Result())

	// the user is subscribed to the newChannelName too
	newChannelFound := false

	for _, subscription := range respGetSubscribedChannels.Subscriptions {
		if subscription.Name == newChannelName {
			newChannelFound = true
			break
		}
	}

	assert.True(t, newChannelFound, "User should be subscribed to the "+newChannelName+" channel")

	// User A subscribes to a new channel
	respSubscribeToChannelUserA, err = userAChannels.SubscribeToChannel(ctx, []channels.SubscribeTo{{Name: "unsubscriber"}})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respSubscribeToChannelUserA.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respSubscribeToChannelUserA.Result())

	// User A unsubscribes from the new channel
	respUnsubscribeFromChannelUserA, err := userAChannels.UnsubscribeFromChannel(ctx, []string{"unsubscriber"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUnsubscribeFromChannelUserA.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respUnsubscribeFromChannelUserA.Result())

	// User A validates the unsubscription
	respGetSubscribedChannelsUserA, err := userAChannels.GetSubscribedChannels(ctx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetSubscribedChannelsUserA.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetSubscribedChannelsUserA.Result())

	unsubscriberChannelFound := false

	for _, subscription := range respGetSubscribedChannelsUserA.Subscriptions {
		if subscription.Name == "unsubscriber" {
			unsubscriberChannelFound = true
			break
		}
	}

	assert.False(t, unsubscriberChannelFound, "User should not be subscribed to the unsubscriber channel")

	// 3 People private chat: Admin, User A, User B
	// Admin sends a message to User A and User B
	respCreatePrivateChat, err := adminMsgSvc.SendMessage(ctx, recipient.ToUsers([]string{userAEmail, userBEmail}), "Hello User A and User B!")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCreatePrivateChat.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respCreatePrivateChat.Result())
	// User A sends a message to Admin and User B
	respCreatePrivateChat, err = userAMsgSvc.SendMessage(ctx, recipient.ToUsers([]string{zulipEmail, userBEmail}), "Hello Admin and User B!")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCreatePrivateChat.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respCreatePrivateChat.Result())
	// User B sends a message to Admin and User A
	userBMsgSvc := messages.NewService(userB)
	respCreatePrivateChat, err = userBMsgSvc.SendMessage(ctx, recipient.ToUsers([]string{zulipEmail, userAEmail}), "Hello Admin and User A!")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCreatePrivateChat.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respCreatePrivateChat.Result())

	// Get Message receipts from the message sent by userB
	respGetMessageReceipts, err := userBMsgSvc.GetMessagesReadReceipts(ctx, respCreatePrivateChat.ID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetMessageReceipts.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetMessageReceipts.Result())
	// Flaky test, probably no one will have read the message
	assert.Empty(t, respGetMessageReceipts.UserIDs)

	// User A marks the message as read
	respMarkAsRead, err := userAMsgSvc.UpdatePersonalMessageFlags(ctx, []int{respCreatePrivateChat.ID}, messages.OperationAdd, messages.FlagRead)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respMarkAsRead.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respMarkAsRead.Result())

	// User B gets the receipts again and should find User A's ID
	respGetMessageReceipts, err = userBMsgSvc.GetMessagesReadReceipts(ctx, respCreatePrivateChat.ID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetMessageReceipts.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetMessageReceipts.Result())
	assert.GreaterOrEqual(t, len(respGetMessageReceipts.UserIDs), 1) // one read receipt
	assert.Contains(t, respGetMessageReceipts.UserIDs, userAID)      // User A has read the message

	// Admin marks the message as read but applying a narrow
	markAsReadNarrow := narrow.NewFilter().
		Add(narrow.New(narrow.ID, respCreatePrivateChat.ID)).              // the message ID
		Add(narrow.New(narrow.Operator("is"), narrow.Operand("private"))). // private messages
		Add(narrow.New(narrow.DmIncluding, narrow.Operand(userAID)))       // including User A in the conversation

	respMarkAsReadNarrow, err := adminMsgSvc.UpdatePersonalMessageFlagsNarrow(ctx, "newest", 1, 1,
		markAsReadNarrow, messages.OperationAdd, messages.FlagRead)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respMarkAsReadNarrow.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respMarkAsReadNarrow.Result())

	// User B gets the receipts again and should find Admin's ID
	respGetMessageReceipts, err = userBMsgSvc.GetMessagesReadReceipts(ctx, respCreatePrivateChat.ID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetMessageReceipts.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetMessageReceipts.Result())
	assert.GreaterOrEqual(t, len(respGetMessageReceipts.UserIDs), 2) // two read receipts
	assert.Contains(t, respGetMessageReceipts.UserIDs, userAID)      // User A has read the message

	// UserA gets its own information
	userAUserSvc := users.NewService(userA)
	respGetUserMe, err := userAUserSvc.GetUserMe(ctx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetUserMe.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetUserMe.Result())
	assert.Equal(t, respGetUserMe.Email, userAEmail)

	// Get User A information by ID
	respGetUser, err := adminUserSvc.GetUser(ctx, respGetUserMe.UserID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetUser.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetUser.Result())
	assert.Equal(t, userAEmail, respGetUser.User.Email)

	// Update User A
	respUpdateUser, err := adminUserSvc.UpdateUser(ctx, respGetUserMe.UserID,
		users.FullName("User A Updated"),
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUpdateUser.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respUpdateUser.Result())
	assert.Empty(t, respUpdateUser.Msg())

	// Get ALL Users
	respGetUsers, err := adminUserSvc.GetUsers(ctx,
		users.ClientGravatars(true),
		users.IncludeCustomProfilesFields(true),
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetUsers.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetUsers.Result())
	assert.GreaterOrEqual(t, len(respGetUsers.Members), 3)

	// Status Text and Emoji
	// Set status text and emoji
	userAStatusText := "I'm busy"
	respUpdateStatus, err := userAUserSvc.UpdateStatus(ctx,
		users.StatusText(userAStatusText),
		users.StatusEmojiName("thumbs_up"),
		// zulip.StatusEmojiCode("1f389"),
		// zulip.StatusReactionType(zulip.ReactionTypeUnicode),
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUpdateStatus.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respUpdateStatus.Result())

	// Get User A status
	respGetUserStatus, err := adminUserSvc.GetUserStatus(ctx, userAID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetUserStatus.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetUserStatus.Result())
	assert.Equal(t, respGetUserStatus.Status.StatusText, userAStatusText)

	// Upload a custom emoji and set it as the user's status emoji
	adminOrgSvc := org.NewService(adminClient)
	emojiName := "dancing_gopher_" + userSuffix
	respUploadEmoji, err := adminOrgSvc.UploadCustomEmoji(ctx, emojiName, "../testdata/dancing_gopher.gif")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUploadEmoji.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respUploadEmoji.Result())

	// Set the custom emoji as the user's status emoji
	respUpdateStatus, err = userAUserSvc.UpdateStatus(ctx,
		users.StatusText("I'm dancing"),
		users.StatusEmojiName(emojiName),
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUpdateStatus.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respUpdateStatus.Result())

	// User B updates his presence to active
	userBUserSvc := users.NewService(userB)
	respUpdatePresence, err := userBUserSvc.UpdateUserPresence(ctx, users.UserPresenceActive)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUpdatePresence.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respUpdatePresence.Result())

	// User A gets presence from User B
	respGetUserPresence, err := userAUserSvc.GetUserPresence(ctx, userBEmail)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetUserPresence.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetUserPresence.Result())
	assert.Equal(t, "active", respGetUserPresence.Presence.Aggregated.Status)

	// Admin Get all users presence and checks User B presence in the list
	respGetUserPresenceAll, err := adminUserSvc.GetUserPresenceAll(ctx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetUserPresenceAll.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetUserPresenceAll.Result())
	assert.Contains(t, respGetUserPresenceAll.Presences, userBEmail)
	assert.Equal(t, "active", respGetUserPresenceAll.Presences[userBEmail].Aggregated.Status)

	// Disable Presence for User B
	respUpdateSettingsPresence, err := userBUserSvc.UpdateSettings(ctx,
		users.PresenceEnabled(false),
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUpdateSettingsPresence.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respUpdateSettingsPresence.Result())

	// Admin gets User B presence and checks that the presence is disabled
	respGetUserBPresence, err := adminUserSvc.GetUserPresence(ctx, userBEmail)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetUserBPresence.HTTPCode())
	assert.Equal(t, zulip.ResultSuccess, respGetUserBPresence.Result())
	assert.Equal(t, "offline", respGetUserBPresence.Presence.Aggregated.Status)

	// spew.Dump("ok")
}
