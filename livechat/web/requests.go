package web

import (
	"context"
	"fmt"
)

func (c *livechatClient) CreateBot(ctx context.Context, payload *CreateBotRequest) (*CreateBotResponse, error) {
	var body CreateBotResponse
	_, err := c.sendRequest(ctx, createBotEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("create_bot action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) DeleteBot(ctx context.Context, payload *DeleteBotRequest) (*DeleteBotResponse, error) {
	var body DeleteBotResponse
	_, err := c.sendRequest(ctx, deleteBotEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("delete_bot action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) ListBots(ctx context.Context, payload *ListBotsRequest) ([]*ListBotResponse, error) {
	var body []*ListBotResponse
	_, err := c.sendRequest(ctx, listBotsEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("list_bots action: %w", err)
	}

	return body, nil
}

func (c *livechatClient) RegisterWebhook(ctx context.Context, payload *RegisterWebhookRequest) (*RegisterWebhookResponse, error) {
	var body RegisterWebhookResponse
	_, err := c.sendRequest(ctx, registerWebhookEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("register_webhook action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) UnregisterWebhook(ctx context.Context, payload *UnregisterWebhookRequest) (*UnregisterWebhookResponse, error) {
	var body UnregisterWebhookResponse
	_, err := c.sendRequest(ctx, registerWebhookEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("unregister_webhook action: %w", err)
	}

	return &body, nil
}
