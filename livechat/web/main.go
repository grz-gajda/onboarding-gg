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

	registerWebhookEndpoint   = "/configuration/action/register_webhook"
	unregisterWebhookEndpoint = "/configuration/action/unregister_webhook"
)

//go:generate mockery --name Client
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type LivechatRequests interface {
	CreateBot(context.Context, *CreateBotRequest) (*CreateBotResponse, error)
	DeleteBot(context.Context, *DeleteBotRequest) (*DeleteBotResponse, error)
	ListBots(context.Context, *ListBotsRequest) ([]*ListBotResponse, error)

	RegisterWebhook(context.Context, *RegisterWebhookRequest) (*RegisterWebhookResponse, error)
	UnregisterWebhook(context.Context, *UnregisterWebhookRequest) (*UnregisterWebhookResponse, error)
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
