package web

import (
	"context"
	"net/http"

	"github.com/livechat/onboarding/livechat"
)

const (
	createBotEndpoint = "/configuration/action/create_bot"
	deleteBotEndpoint = "/configuration/action/delete_bot"
	listBotsEndpoint  = "/configuration/action/list_bots"

	transferChatEndpoint = "/agent/action/transfer_chat"
	sendEventEndpoint    = "/agent/action/send_event"

	registerWebhookEndpoint       = "/configuration/action/register_webhook"
	unregisterWebhookEndpoint     = "/configuration/action/unregister_webhook"
	enableLicenseWebhookEndpoint  = "/configuration/action/enable_license_webhooks"
	disableLicenseWebhookEndpoint = "/configuration/action/disable_license_webhooks"

	setRoutingStatusEndpoint = "/agent/action/set_routing_status"
)

//go:generate mockery --name Client
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type LivechatRequests interface {
	CreateBot(context.Context, *CreateBotRequest) (*CreateBotResponse, error)
	DeleteBot(context.Context, *DeleteBotRequest) (*DeleteBotResponse, error)
	ListBots(context.Context, *ListBotsRequest) ([]*ListBotResponse, error)

	TransferChat(context.Context, *TransferChatRequest) (*TransferChatResponse, error)
	SendEvent(context.Context, *SendEventRequest) (*SendEventResponse, error)

	RegisterWebhook(context.Context, *RegisterWebhookRequest) (*RegisterWebhookResponse, error)
	UnregisterWebhook(context.Context, *UnregisterWebhookRequest) (*UnregisterWebhookResponse, error)
	EnableLicenseWebhook(context.Context, *EnableLicenseWebhookRequest) (*EnableLicenseWebhookResponse, error)
	DisableLicenseWebhook(context.Context, *DisableLicenseWebhookRequest) (*DisableLicenseWebhookResponse, error)

	SetRoutingStatus(context.Context, *SetRoutingStatusRequest) (*SetRoutingStatusResponse, error)
}

func New(client Client, url string) LivechatRequests {
	return &livechatClient{
		httpClient: client,
		url:        url,
	}
}

type CreateBotRequest struct {
	Name     string            `json:"name"`
	ClientID livechat.ClientID `json:"owner_client_id,omitempty"`
}

type CreateBotResponse struct {
	ID livechat.AgentID `json:"id"`
}

type DeleteBotRequest struct {
	ID livechat.AgentID `json:"id"`
}

type DeleteBotResponse struct{}

type ListBotsRequest struct {
	All bool `json:"all,omitempty"`
}

type ListBotResponse struct {
	ID   livechat.AgentID `json:"id"`
	Name string           `json:"name"`
}

type RegisterWebhookRequest struct {
	Action    string            `json:"action"`
	SecretKey string            `json:"secret_key"`
	URL       string            `json:"url"`
	ClientID  livechat.ClientID `json:"owner_client_id,omitempty"`
	Type      string            `json:"type"`
}

type RegisterWebhookResponse struct {
	ID string `json:"id"`
}

type UnregisterWebhookRequest struct {
	ID       string            `json:"id"`
	ClientID livechat.ClientID `json:"owner_client_id,omitempty"`
}

type UnregisterWebhookResponse struct{}

type TransferChatRequest struct {
	ID     livechat.ChatID `json:"id"`
	Target struct {
		Type string             `json:"type"`
		IDs  []livechat.AgentID `json:"ids"`
	} `json:"target"`
	Force bool `json:"force,omitempty"`
}

type TransferChatResponse struct{}

type Event struct {
	Text       string `json:"text"`
	Type       string `json:"type"`
	Recipients string `json:"recipients"`
}

type SendEventRequest struct {
	ChatID livechat.ChatID `json:"chat_id"`
	Event  Event           `json:"event"`
}

type SendEventResponse struct {
	EventID string `json:"event_id"`
}

type EnableLicenseWebhookRequest struct {
	ClientID livechat.ClientID `json:"owner_client_id,omitempty"`
}

type EnableLicenseWebhookResponse struct{}

type DisableLicenseWebhookRequest struct {
	ClientID livechat.ClientID `json:"owner_client_id,omitempty"`
}

type DisableLicenseWebhookResponse struct{}

type SetRoutingStatusRequest struct {
	Status  string           `json:"status"`
	AgentID livechat.AgentID `json:"agent_id"`
}

type SetRoutingStatusResponse struct{}
