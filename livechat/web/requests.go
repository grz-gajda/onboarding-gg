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

func (c *livechatClient) TransferChat(ctx context.Context, payload *TransferChatRequest) (*TransferChatResponse, error) {
	var body TransferChatResponse
	_, err := c.sendRequest(ctx, transferChatEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("transfer_chat action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) SendEvent(ctx context.Context, payload *SendEventRequest) (*SendEventResponse, error) {
	var body SendEventResponse
	_, err := c.sendRequest(ctx, sendEventEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("send_event action: %w", err)
	}

	return &body, nil
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
	_, err := c.sendRequest(ctx, unregisterWebhookEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("unregister_webhook action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) EnableLicenseWebhook(ctx context.Context, payload *EnableLicenseWebhookRequest) (*EnableLicenseWebhookResponse, error) {
	var body EnableLicenseWebhookResponse
	_, err := c.sendRequest(ctx, enableLicenseWebhookEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("enable_license_webhook action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) DisableLicenseWebhook(ctx context.Context, payload *DisableLicenseWebhookRequest) (*DisableLicenseWebhookResponse, error) {
	var body DisableLicenseWebhookResponse
	_, err := c.sendRequest(ctx, disableLicenseWebhookEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("disable_license_webhook action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) SetRoutingStatus(ctx context.Context, payload *SetRoutingStatusRequest) (*SetRoutingStatusResponse, error) {
	var body SetRoutingStatusResponse
	_, err := c.sendRequest(ctx, setRoutingStatusEndpoint, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("set_routing_status action: %w", err)
	}

	return &body, nil
}
