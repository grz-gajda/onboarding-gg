package bot_webhooks

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
	log "github.com/sirupsen/logrus"
)

type app struct {
	lcHTTP    web.LivechatRequests
	licenseID livechat.LicenseID
	agents    *agents
	webhooks  map[string]*webhookDetails
}

type webhookDetails struct {
	id string
}

func newApp(lcHTTP web.LivechatRequests, id livechat.LicenseID) *app {
	return &app{
		lcHTTP:    lcHTTP,
		licenseID: id,
		agents:    &agents{},
		webhooks:  make(map[string]*webhookDetails),
	}
}

func (a *app) RegisterAction(ctx context.Context, action string, opts *registerActionOptions) error {
	clientID, err := auth.GetClientID(ctx)
	if err != nil {
		return fmt.Errorf("bot: register_action: %w", err)
	}

	payload := &web.RegisterWebhookRequest{
		ClientID:  clientID,
		SecretKey: "random secret key",
		URL:       fmt.Sprintf("%s/%s", "http://localhost:8081/webhooks", action),
		Action:    action,
		Type:      "bot",
	}

	webhookResponse, err := a.lcHTTP.RegisterWebhook(ctx, payload)
	if err != nil {
		return fmt.Errorf("bot: register_action: %w", err)
	}

	a.webhooks[action] = &webhookDetails{id: webhookResponse.ID}
	return nil
}

func (a *app) UnregisterActions(ctx context.Context) error {
	clientID, err := auth.GetClientID(ctx)
	if err != nil {
		return fmt.Errorf("bot: register_action: %w", err)
	}

	wg := sync.WaitGroup{}

	for actionName, details := range a.webhooks {
		wg.Add(1)

		go func(aName string, d *webhookDetails) {
			defer wg.Done()
			_, err := a.lcHTTP.UnregisterWebhook(ctx, &web.UnregisterWebhookRequest{
				ID:       d.id,
				ClientID: clientID,
			})

			logEntry := log.WithField("license_id", a.licenseID).WithField("webhook_id", d.id).WithContext(ctx).WithField("action", aName)
			if err != nil {
				logEntry.WithError(err).Error("Cannot unregister webhook")
			} else {
				logEntry.Debug("Webhook unregistered")
			}
		}(actionName, details)
	}

	wg.Wait()
	return nil
}

func (a *app) TransferChat(ctx context.Context, msg *rtm.PushIncomingChat) error {
	log.WithField("agents", a.agents).Debug("TransferChat action")
	agent, err := a.agents.FindByChatExclude(msg.Payload.Chat.ID)
	if err != nil {
		return fmt.Errorf("bot: transfer_chat action: %w", err)
	}

	_, err = a.lcHTTP.TransferChat(ctx, &web.TransferChatRequest{
		ID: msg.Payload.Chat.ID,
		Target: struct {
			Type string             "json:\"type\""
			IDs  []livechat.AgentID "json:\"ids\""
		}{
			Type: "agent",
			IDs:  []livechat.AgentID{agent.ID},
		},
	})

	if err != nil {
		if !strings.Contains(err.Error(), "One or more of requested agents are already present in the chat") {
			return fmt.Errorf("bot: transfer_chat action: %w", err)
		}
	}

	agent.chats = append(agent.chats, msg.Payload.Chat.ID)
	return nil
}

func (a *app) CreateBot(ctx context.Context) error {
	clientID, err := auth.GetClientID(ctx)
	if err != nil {
		return fmt.Errorf("bot: create bot: %w", err)
	}

	response, err := a.lcHTTP.CreateBot(ctx, &web.CreateBotRequest{
		Name:     "OnboardingGG",
		ClientID: clientID,
	})
	if err != nil {
		return fmt.Errorf("bot: create_bot: %w", err)
	}

	agent := newAgent(response.ID, a.lcHTTP)
	if err := a.agents.Register(agent); err != nil {
		return fmt.Errorf("bot: create_bot: %w", err)
	}

	log.WithField("license_id", a.licenseID).WithField("amount", 1).Info("Bot has been registered")
	return nil
}

func (a *app) FetchBots(ctx context.Context) error {
	botsResponse, err := a.lcHTTP.ListBots(ctx, &web.ListBotsRequest{All: true})
	if err != nil {
		return fmt.Errorf("bot: fetch_bots: %w", err)
	}
	if len(botsResponse) == 0 {
		if err := a.CreateBot(ctx); err != nil {
			return fmt.Errorf("bot: fetch_bots: %w", err)
		}
		return nil
	}

	for _, botID := range botsResponse {
		agent := newAgent(botID.ID, a.lcHTTP)
		if err := a.agents.Register(agent); err != nil {
			return fmt.Errorf("bot: fetch_bots: %w", err)
		}
	}

	log.WithField("license_id", a.licenseID).WithField("amount", len(botsResponse)).Infof("Bot has been registered")
	return nil
}

func (a *app) IncomingEvent(ctx context.Context, msg *rtm.PushIncomingMessage) error {
	log.WithField("agents", a.agents).WithField("chat_id", msg.Payload.ChatID).Debug("IncomingEvent action")

	agent, err := a.agents.FindByChat(msg.Payload.ChatID)
	if err != nil {
		return fmt.Errorf("bot: incoming_event action: %w", err)
	}

	ctx = auth.WithAuthorID(ctx, agent.ID)

	_, err = a.lcHTTP.SendEvent(ctx, &web.SendEventRequest{
		ChatID: msg.Payload.ChatID,
		Event: web.Event{
			Text:       "Lorem ipsum dolor sit amet",
			Type:       "message",
			Recipients: "all",
		},
	})

	return fmt.Errorf("bot: incoming_event action: %w", err)
}

type registerActionOptions struct {
	CustomURL string
}
