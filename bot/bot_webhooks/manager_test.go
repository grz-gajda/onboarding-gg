package bot_webhooks

import (
	"context"
	"fmt"
	"testing"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/web/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	validLicenseID   = livechat.LicenseID(23456)
	invalidLicenseID = livechat.LicenseID(12345)

	validChatID = livechat.ChatID("chat_id_1")
	validBotID  = livechat.AgentID("bot_1")
)

var (
	webhooksLen = len(WebhookEvents)
)

func Test_Manager_Install(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	manager, err := helperCreateManager(t, ctx, lcHTTP)
	assert.NoError(t, err)
	assert.Len(t, manager.apps.apps, 1)

	lcHTTP.AssertNumberOfCalls(t, "RegisterWebhook", webhooksLen)
	lcHTTP.AssertNumberOfCalls(t, "EnableLicenseWebhook", 1)

	app := manager.apps.apps[0]
	assert.NotNil(t, app)
	assert.Equal(t, 1, app.agents.Len())

	a, unlock := app.agents.Get()
	defer unlock()

	assert.Equal(t, validBotID, a[0].ID)
}

func Test_Manager_Uninstall_InvalidLicenseID(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	manager, err := helperCreateManager(t, ctx, lcHTTP)
	assert.NoError(t, err)
	assert.Error(t, manager.UninstallApp(ctx, invalidLicenseID))
}

func Test_Manager_Uninstall_ValidLicenseID(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	// +uninstall webhooks
	lcHTTP.On("DisableLicenseWebhook", ctx, mock.Anything).Once().Return(&livechat.DisableLicenseWebhookResponse{}, nil)
	lcHTTP.On("UnregisterWebhook", ctx, mock.Anything).Times(webhooksLen).Return(&livechat.UnregisterWebhookResponse{}, nil)
	// +uninstall bots
	lcHTTP.On("SetRoutingStatus", ctx, mock.Anything).Twice().Return(&livechat.SetRoutingStatusResponse{}, nil)

	manager, err := helperCreateManager(t, ctx, lcHTTP)
	assert.NoError(t, err)
	assert.NoError(t, manager.UninstallApp(ctx, validLicenseID))
}

func Test_Manager_Redirect_IncomingChat(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	lcHTTP.On("TransferChat", ctx, mock.MatchedBy(func(p *livechat.TransferChatRequest) bool {
		return p.ID == validChatID
	})).Once().Return(&livechat.TransferChatResponse{}, nil)

	manager, _ := helperCreateManager(t, ctx, lcHTTP)
	err := manager.Redirect(ctx, helperBuildPushIncomingChat(t, validLicenseID, validChatID))
	assert.NoError(t, err)
}

func Test_Manager_Redirect_IncomingEvent(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	lcHTTP.On("SendEvent", mock.MatchedBy(func(c context.Context) bool {
		authorID, err := auth.GetAuthorID(c)
		return assert.NoError(t, err) && assert.Equal(t, validBotID, authorID)
	}), mock.MatchedBy(func(p *livechat.Event) bool {
		return p.ChatID == validChatID
	})).Once().Return(&livechat.SendEventResponse{}, nil)

	message := helperBuildPushIncomingEvent(t, validLicenseID, validChatID)
	message.Payload.Event.Text = "Hello world"
	message.Payload.Event.AuthorID = "custom_author_id"

	manager, _ := helperCreateManager(t, ctx, lcHTTP)
	err := manager.Redirect(ctx, message)
	assert.NoError(t, err)
}

func Test_Manager_Redirect_IncomingEvent_TransferChat(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	lcHTTP.On("TransferChat", mock.Anything, mock.MatchedBy(func(p *livechat.TransferChatRequest) bool {
		return p.ID == validChatID && len(p.Target.IDs) > 0
	})).Twice().Return(&livechat.TransferChatResponse{}, nil)

	lcHTTP.On("ListAgents", mock.Anything, mock.Anything).Once().Return([]*livechat.ListAgentsResponse{
		{ID: livechat.AgentID("agent_1234")},
	}, nil)

	message := helperBuildPushIncomingEvent(t, validLicenseID, validChatID)
	message.Payload.Event.Text = "Wróć do człowieka"
	message.Payload.Event.AuthorID = "custom_author_id"

	manager, _ := helperCreateManager(t, ctx, lcHTTP)
	assert.NoError(t, manager.Redirect(ctx, helperBuildPushIncomingChat(t, validLicenseID, validChatID)))
	assert.NoError(t, manager.Redirect(ctx, message))
	lcHTTP.AssertNumberOfCalls(t, "TransferChat", 2)
}

func Test_Manager_UserAddedToChat(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	lcHTTP.On("TransferChat", mock.Anything, mock.Anything).Once().Return(&livechat.TransferChatResponse{}, nil)

	manager, _ := helperCreateManager(t, ctx, lcHTTP)
	assert.NoError(t, manager.Redirect(ctx, helperBuildPushIncomingChat(t, validLicenseID, validChatID)))
	assert.Equal(t, 1, manager.apps.apps[0].agents.Len())

	assert.NoError(t, manager.Redirect(ctx, helperBuildPushUserAddedToChat(t, validLicenseID, validChatID)))
	_, err := manager.apps.apps[0].agents.FindByChat(validChatID)
	assert.Error(t, err)
}

func helperCreateManager(t *testing.T, ctx context.Context, lcHTTP *mocks.LivechatRequests) (*manager, error) {
	t.Helper()
	// +install bots
	lcHTTP.On("CreateBot", ctx, mock.Anything).Once().Return(&livechat.CreateBotResponse{ID: validBotID}, nil)
	lcHTTP.On("ListBots", ctx, mock.Anything).Once().Return([]*livechat.ListBotResponse{{ID: validBotID}}, nil)
	lcHTTP.On("SetRoutingStatus", ctx, mock.Anything).Once().Return(&livechat.SetRoutingStatusResponse{}, nil)
	// +install webhooks
	lcHTTP.On("RegisterWebhook", ctx, mock.Anything).Times(webhooksLen).Return(&livechat.RegisterWebhookResponse{}, nil)
	lcHTTP.On("EnableLicenseWebhook", ctx, mock.Anything).Twice().Return(&livechat.EnableLicenseWebhookResponse{}, nil)

	mng := New(lcHTTP, "http://localhost:8081", "author_id")
	err := mng.InstallApp(ctx, validLicenseID)

	if !assert.NoError(t, err) {
		return nil, fmt.Errorf("canno create manager: %w", err)
	}

	rawManager, ok := mng.(*manager)
	if !ok {
		return nil, fmt.Errorf("cannot create manager: %w", err)
	}

	return rawManager, nil
}

func helperBuildPushIncomingChat(t *testing.T, licenseID livechat.LicenseID, chatID livechat.ChatID) *livechat.PushIncomingChat {
	t.Helper()
	return &livechat.PushIncomingChat{
		Action:    "incoming_chat",
		LicenseID: licenseID,
		Payload: struct {
			Chat struct {
				ID livechat.ChatID "json:\"id\""
			} "json:\"chat\""
		}{
			Chat: struct {
				ID livechat.ChatID "json:\"id\""
			}{
				ID: chatID,
			},
		},
	}
}

func helperBuildPushIncomingEvent(t *testing.T, licenseID livechat.LicenseID, chatID livechat.ChatID) *livechat.PushIncomingMessage {
	t.Helper()

	return &livechat.PushIncomingMessage{
		Action:    "incoming_message",
		LicenseID: licenseID,
		Payload: struct {
			ChatID   livechat.ChatID "json:\"chat_id\""
			ThreadID string          "json:\"thread_id,omitempty\""
			Event    struct {
				Type     string "json:\"type\""
				Text     string "json:\"text\""
				AuthorID string "json:\"author_id\""
			} "json:\"event\""
		}{
			ChatID: chatID,
			Event: struct {
				Type     string "json:\"type\""
				Text     string "json:\"text\""
				AuthorID string "json:\"author_id\""
			}{
				Type: "message",
			},
		},
	}
}

func helperBuildPushUserAddedToChat(t *testing.T, licenseID livechat.LicenseID, chatID livechat.ChatID) *livechat.PushUserAddedToChat {
	t.Helper()
	return &livechat.PushUserAddedToChat{
		Action:    "user_added_to_chat",
		LicenseID: licenseID,
		Payload: struct {
			ChatID livechat.ChatID "json:\"chat_id\""
			User   struct {
				Present bool   "json:\"present\""
				Type    string "json:\"type\""
			} "json:\"user\""
		}{
			ChatID: chatID,
			User: struct {
				Present bool   "json:\"present\""
				Type    string "json:\"type\""
			}{
				Present: true,
				Type:    "agent",
			},
		},
	}
}
