package web

import (
	"context"

	"github.com/livechat/onboarding/livechat"
)

//go:generate mockery --name LivechatRequests
type LivechatRequests interface {
	CreateBot(context.Context, *livechat.CreateBotRequest) (*livechat.CreateBotResponse, error)
	DeleteBot(context.Context, *livechat.DeleteBotRequest) (*livechat.DeleteBotResponse, error)
	ListBots(context.Context, *livechat.ListBotsRequest) ([]*livechat.ListBotResponse, error)
	ListAgents(context.Context, *livechat.ListAgentsRequest) ([]*livechat.ListAgentsResponse, error)

	TransferChat(context.Context, *livechat.TransferChatRequest) (*livechat.TransferChatResponse, error)
	SendEvent(context.Context, *livechat.Event) (*livechat.SendEventResponse, error)

	RegisterWebhook(context.Context, *livechat.RegisterWebhookRequest) (*livechat.RegisterWebhookResponse, error)
	UnregisterWebhook(context.Context, *livechat.UnregisterWebhookRequest) (*livechat.UnregisterWebhookResponse, error)
	EnableLicenseWebhook(context.Context, *livechat.EnableLicenseWebhookRequest) (*livechat.EnableLicenseWebhookResponse, error)
	DisableLicenseWebhook(context.Context, *livechat.DisableLicenseWebhookRequest) (*livechat.DisableLicenseWebhookResponse, error)

	SetRoutingStatus(context.Context, *livechat.SetRoutingStatusRequest) (*livechat.SetRoutingStatusResponse, error)
}

func New(client livechat.Client, url string) LivechatRequests {
	return &livechatClient{
		httpClient: client,
		url:        url,
	}
}
