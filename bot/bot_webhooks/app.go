package bot_webhooks

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/livechat/onboarding/bot"
	"github.com/livechat/onboarding/bot/bot_webhooks/agents"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/web"
	log "github.com/sirupsen/logrus"
)

type app struct {
	lcHTTP    web.LivechatRequests
	sender    bot.Sender
	licenseID livechat.LicenseID
	agents    agents.Agents
	webhooks  map[string]*webhookDetails
	localURL  string
}

type webhookDetails struct {
	id string
}

func newApp(lcHTTP web.LivechatRequests, sender bot.Sender, id livechat.LicenseID, localURL string) *app {
	return &app{
		lcHTTP:    lcHTTP,
		licenseID: id,
		agents:    agents.NewCollection(),
		webhooks:  make(map[string]*webhookDetails),
		localURL:  localURL,
		sender:    sender,
	}
}

func (a *app) RegisterAction(ctx context.Context, actions ...string) error {
	for _, action := range actions {
		payload := &livechat.RegisterWebhookRequest{
			SecretKey: "random secret key",
			URL:       fmt.Sprintf("%s/webhooks/%s", a.localURL, action),
			Action:    action,
			Type:      "license",
		}

		webhookResponse, err := a.lcHTTP.RegisterWebhook(ctx, payload)
		if err != nil {
			log.WithField("action", action).WithError(err).Error("Something went wrong with webhook's registration")
			return fmt.Errorf("bot: register_action: %w", err)
		}

		log.WithField("action", action).WithField("url", fmt.Sprintf("%s/webhooks/%s", a.localURL, action)).Debug("Webhook registered")

		a.webhooks[action] = &webhookDetails{id: webhookResponse.ID}
	}

	return nil
}

func (a *app) UnregisterActions(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		agents.Terminate(ctx, a.lcHTTP, a.agents)
	}()

	for actionName, details := range a.webhooks {
		wg.Add(1)

		go func(aName string, d *webhookDetails) {
			defer wg.Done()
			_, err := a.lcHTTP.UnregisterWebhook(ctx, &livechat.UnregisterWebhookRequest{ID: d.id})
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

func (a *app) TransferChat(ctx context.Context, msg *livechat.PushIncomingChat) error {
	agent, err := a.agents.FindByChatExclude(msg.Payload.Chat.ID)
	if err != nil {
		return fmt.Errorf("bot: transfer_chat action: %w", err)
	}

	if _, err = a.lcHTTP.TransferChat(ctx, buildTransferChatMessage(msg.Payload.Chat.ID, agent.ID)); !isTransferChatErrorOk(err) {
		return fmt.Errorf("bot: transfer_chat action: %w", err)
	}

	return agent.RegisterChat(msg.Payload.Chat.ID)
}

func (a *app) IncomingEvent(ctx context.Context, msg *livechat.PushIncomingMessage) error {
	agent, err := a.agents.FindByChat(msg.Payload.ChatID)
	if err != nil {
		if err := a.TransferChat(ctx, mapIncomingEventIntoTransferChat(msg)); err != nil {
			return fmt.Errorf("bot: incoming_event action: %w", err)
		}
		return a.IncomingEvent(ctx, msg)
	}

	return a.sender.Talk(auth.WithAuthorID(ctx, agent.ID), msg.Payload.ChatID, msg)
}

func buildTransferChatMessage(chatID livechat.ChatID, agentID livechat.AgentID) *livechat.TransferChatRequest {
	return &livechat.TransferChatRequest{
		ID: chatID,
		Target: struct {
			Type string             "json:\"type\""
			IDs  []livechat.AgentID "json:\"ids\""
		}{
			Type: "agent",
			IDs:  []livechat.AgentID{agentID},
		},
		Force: true,
	}
}

func mapIncomingEventIntoTransferChat(msg *livechat.PushIncomingMessage) *livechat.PushIncomingChat {
	return &livechat.PushIncomingChat{
		Action:    "incoming_chat",
		LicenseID: msg.LicenseID,
		Payload: struct {
			Chat struct {
				ID livechat.ChatID "json:\"id\""
			} "json:\"chat\""
		}{
			Chat: struct {
				ID livechat.ChatID "json:\"id\""
			}{
				ID: msg.Payload.ChatID,
			},
		},
	}
}

func isTransferChatErrorOk(err error) bool {
	return err == nil || strings.Contains(err.Error(), "One or more of requested agents are already present in the chat")
}
