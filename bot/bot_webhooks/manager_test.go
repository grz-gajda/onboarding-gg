package bot_webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	lcMocks "github.com/livechat/onboarding/livechat/mocks"
	"github.com/livechat/onboarding/livechat/web/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	validLicenseID   = livechat.LicenseID(23456)
	invalidLicenseID = livechat.LicenseID(12345)

	validChatID = livechat.ChatID("chat_id_1")
	validBotID  = livechat.AgentID("bot_1")

	oauthToken = "some_random_token"
)

var (
	webhooksLen = len(WebhookEvents)
	matchCtx    = mock.MatchedBy(mockContextWithOAuthToken)
)

func Test_Manager_Authorize(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	httpClient := new(lcMocks.Client)
	httpClient.On("Do", mock.Anything).Return(&http.Response{
		Body:       io.NopCloser(strings.NewReader(`{"access_token":"abcd"}`)),
		StatusCode: http.StatusOK,
	}, nil).Once()
	httpClient.On("Do", mock.Anything).Return(&http.Response{
		Body:       io.NopCloser(strings.NewReader(`{"access_token":"abcd"}`)),
		StatusCode: http.StatusOK,
	}, nil).Once()

	mng := New(lcHTTP, "", "")
	assert.NoError(t, mng.Authorize(ctx, httpClient, &auth.AuthorizeCredentials{}))
	assert.NoError(t, mng.Authorize(ctx, httpClient, &auth.AuthorizeCredentials{}))
	httpClient.AssertNumberOfCalls(t, "Do", 2)
}

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
	lcHTTP.On("DisableLicenseWebhook", matchCtx, mock.Anything).Once().Return(&livechat.DisableLicenseWebhookResponse{}, nil)
	lcHTTP.On("UnregisterWebhook", matchCtx, mock.Anything).Times(webhooksLen).Return(&livechat.UnregisterWebhookResponse{}, nil)
	// +uninstall bots
	lcHTTP.On("SetRoutingStatus", matchCtx, mock.Anything).Twice().Return(&livechat.SetRoutingStatusResponse{}, nil)

	manager, err := helperCreateManager(t, ctx, lcHTTP)
	assert.NoError(t, err)
	assert.NoError(t, manager.UninstallApp(ctx, validLicenseID))
}

func Test_Manager_Redirect_IncomingChat(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	lcHTTP.On("TransferChat", matchCtx, mock.MatchedBy(func(p *livechat.TransferChatRequest) bool {
		return p.ID == validChatID
	})).Once().Return(&livechat.TransferChatResponse{}, nil)

	manager, _ := helperCreateManager(t, ctx, lcHTTP)
	err := manager.Redirect(ctx, helperBuildPushIncomingChat(t, validLicenseID, validChatID))
	assert.NoError(t, err)
}

func Test_Manager_Redirect_IncomingEvent(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	message := helperBuildPushIncomingEvent(t, validLicenseID, validChatID)
	message.Payload.Event.Text = "Hello world"
	message.Payload.Event.AuthorID = "custom_author_id"

	manager, _ := helperCreateManager(t, ctx, lcHTTP)
	err := manager.Redirect(ctx, message)
	assert.NoError(t, err)
	lcHTTP.AssertNumberOfCalls(t, "SendEvent", 0)
}

func Test_Manager_Redirect_IncomingEvent_TransferChat(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	lcHTTP.On("TransferChat", matchCtx, mock.MatchedBy(func(p *livechat.TransferChatRequest) bool {
		return p.ID == validChatID && len(p.Target.IDs) > 0
	})).Twice().Return(&livechat.TransferChatResponse{}, nil)

	lcHTTP.On("ListAgentsForTransfer", matchCtx, mock.Anything).Once().Return([]*livechat.ListAgentsForTransferResponse{
		{AgentID: livechat.AgentID("agent_1234")},
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

	lcHTTP.On("TransferChat", matchCtx, mock.Anything).Once().Return(&livechat.TransferChatResponse{}, nil)

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
	lcHTTP.On("CreateBot", matchCtx, mock.Anything).Once().Return(&livechat.CreateBotResponse{ID: validBotID}, nil)
	lcHTTP.On("ListBots", matchCtx, mock.Anything).Once().Return([]*livechat.ListBotResponse{{ID: validBotID}}, nil)
	lcHTTP.On("SetRoutingStatus", matchCtx, mock.Anything).Once().Return(&livechat.SetRoutingStatusResponse{}, nil)
	// +install webhooks
	lcHTTP.On("RegisterWebhook", matchCtx, mock.Anything).Times(webhooksLen).Return(&livechat.RegisterWebhookResponse{}, nil)
	lcHTTP.On("EnableLicenseWebhook", matchCtx, mock.Anything).Twice().Return(&livechat.EnableLicenseWebhookResponse{}, nil)

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			cancel()
		}
	}()

	mng := New(lcHTTP, "http://localhost:8081", "author_id")
	go func() {
		byteBody, err := json.Marshal(map[string]string{"access_token": oauthToken})
		if err != nil {
			return
		}

		httpClient := new(lcMocks.Client)
		httpClient.On("Do", mock.Anything).Once().Return(&http.Response{
			Body:       io.NopCloser(bytes.NewBuffer(byteBody)),
			StatusCode: http.StatusOK,
		}, nil)

		time.Sleep(100 * time.Millisecond)
		assert.NoError(t, mng.Authorize(ctx, httpClient, &auth.AuthorizeCredentials{}))
	}()

	err := mng.InstallApp(ctx, validLicenseID)
	if !assert.NoError(t, err) {
		t.Fatalf("cannot create manager: %s", err.Error())
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

func helperBuildGetChatResponse(t *testing.T, chatID livechat.ChatID, agentsID ...livechat.AgentID) *livechat.GetChatResponse {
	users := []struct {
		ID   livechat.AgentID "json:\"id\""
		Type string           "json:\"type\""
	}{}
	for _, agentID := range agentsID {
		users = append(users, struct {
			ID   livechat.AgentID "json:\"id\""
			Type string           "json:\"type\""
		}{
			ID:   agentID,
			Type: "agent",
		})
	}

	return &livechat.GetChatResponse{
		ID:      chatID,
		UserIDs: agentsID,
		Users:   users,
	}
}

func mockContextWithOAuthToken(ctx context.Context) bool {
	token, err := auth.GetAuthToken(ctx)
	if err != nil {
		return false
	}

	return token == fmt.Sprintf("Bearer %s", oauthToken)
}
