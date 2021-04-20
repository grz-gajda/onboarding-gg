package bot

import (
	"context"
	"errors"
	"testing"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/web/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	definedLicenseID = livechat.LicenseID(1234)
	definedChatID    = livechat.ChatID("custom_chat_id")
	definedAuthorID  = "custom_author_id"
)

func Test_Sender_Hello(t *testing.T) {
	lcHTTP := new(mocks.LivechatRequests)
	ctx := context.Background()

	lcHTTP.On("SendEvent", ctx, mock.MatchedBy(func(p *livechat.Event) bool {
		return p.Event.Text == "World!"
	})).Return(&livechat.SendEventResponse{}, nil)

	msg := helperBuildPushIncomingEvent(t, definedLicenseID, definedChatID)
	msg.Payload.Event.Text = "Hello"

	sender := NewSender(lcHTTP, definedAuthorID)
	assert.NoError(t, sender.Talk(ctx, definedChatID, msg))
}

func Test_Sender_Transfer_AgentOffline(t *testing.T) {
	lcHTTP := new(mocks.LivechatRequests)
	ctx := context.Background()

	lcHTTP.On("SendEvent", ctx, mock.MatchedBy(func(p *livechat.Event) bool {
		return p.Event.Text == "Obecnie nie ma żadnego człowieka do rozmowy :("
	})).Return(&livechat.SendEventResponse{}, nil)

	lcHTTP.On("ListAgents", ctx, mock.Anything).Return([]*livechat.ListAgentsResponse{{
		ID: "abcd",
	}}, nil)

	lcHTTP.On("TransferChat", ctx, mock.MatchedBy(func(p *livechat.TransferChatRequest) bool {
		return p.Target.IDs[0] == "abcd"
	})).Return(nil, errors.New("Agent is offline."))

	msg := helperBuildPushIncomingEvent(t, definedLicenseID, definedChatID)
	msg.Payload.Event.Text = "Wróć do człowieka"

	sender := NewSender(lcHTTP, definedAuthorID)
	assert.NoError(t, sender.Talk(ctx, definedChatID, msg))
	lcHTTP.AssertNumberOfCalls(t, "TransferChat", 1)
	lcHTTP.AssertNumberOfCalls(t, "SendEvent", 1)
}

func Test_Sender_Transfer_AgentOnline(t *testing.T) {
	lcHTTP := new(mocks.LivechatRequests)
	ctx := context.Background()

	lcHTTP.On("ListAgents", ctx, mock.Anything).Return([]*livechat.ListAgentsResponse{{
		ID: "abcd",
	}}, nil)

	lcHTTP.On("TransferChat", ctx, mock.Anything).Return(&livechat.TransferChatResponse{}, nil)

	msg := helperBuildPushIncomingEvent(t, definedLicenseID, definedChatID)
	msg.Payload.Event.Text = "Wróć do człowieka"

	sender := NewSender(lcHTTP, definedAuthorID)
	assert.NoError(t, sender.Talk(ctx, definedChatID, msg))
	lcHTTP.AssertNumberOfCalls(t, "TransferChat", 1)
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
