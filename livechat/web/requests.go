package web

import (
	"context"
	"fmt"

	"github.com/livechat/onboarding/livechat"
)

func (c *livechatClient) CreateBot(ctx context.Context, payload *livechat.CreateBotRequest) (*livechat.CreateBotResponse, error) {
	var body livechat.CreateBotResponse
	_, err := c.sendRequest(ctx, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("create_bot action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) DeleteBot(ctx context.Context, payload *livechat.DeleteBotRequest) (*livechat.DeleteBotResponse, error) {
	var body livechat.DeleteBotResponse
	_, err := c.sendRequest(ctx, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("delete_bot action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) ListBots(ctx context.Context, payload *livechat.ListBotsRequest) ([]*livechat.ListBotResponse, error) {
	var body []*livechat.ListBotResponse
	_, err := c.sendRequest(ctx, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("list_bots action: %w", err)
	}

	return body, nil
}

func (c *livechatClient) TransferChat(ctx context.Context, payload *livechat.TransferChatRequest) (*livechat.TransferChatResponse, error) {
	var body livechat.TransferChatResponse
	_, err := c.sendRequest(ctx, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("transfer_chat action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) SendEvent(ctx context.Context, payload *livechat.Event) (*livechat.SendEventResponse, error) {
	var body livechat.SendEventResponse
	_, err := c.sendRequest(ctx, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("send_event action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) RegisterWebhook(ctx context.Context, payload *livechat.RegisterWebhookRequest) (*livechat.RegisterWebhookResponse, error) {
	var body livechat.RegisterWebhookResponse
	_, err := c.sendRequest(ctx, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("register_webhook action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) UnregisterWebhook(ctx context.Context, payload *livechat.UnregisterWebhookRequest) (*livechat.UnregisterWebhookResponse, error) {
	var body livechat.UnregisterWebhookResponse
	_, err := c.sendRequest(ctx, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("unregister_webhook action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) EnableLicenseWebhook(ctx context.Context, payload *livechat.EnableLicenseWebhookRequest) (*livechat.EnableLicenseWebhookResponse, error) {
	var body livechat.EnableLicenseWebhookResponse
	_, err := c.sendRequest(ctx, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("enable_license_webhook action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) DisableLicenseWebhook(ctx context.Context, payload *livechat.DisableLicenseWebhookRequest) (*livechat.DisableLicenseWebhookResponse, error) {
	var body livechat.DisableLicenseWebhookResponse
	_, err := c.sendRequest(ctx, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("disable_license_webhook action: %w", err)
	}

	return &body, nil
}

func (c *livechatClient) SetRoutingStatus(ctx context.Context, payload *livechat.SetRoutingStatusRequest) (*livechat.SetRoutingStatusResponse, error) {
	var body livechat.SetRoutingStatusResponse
	_, err := c.sendRequest(ctx, payload, &body)
	if err != nil {
		return nil, fmt.Errorf("set_routing_status action: %w", err)
	}

	return &body, nil
}
