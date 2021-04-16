package bot_webhooks

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/livechat/onboarding/bot/messages"
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
	localURL  string
}

type webhookDetails struct {
	id string
}

func newApp(lcHTTP web.LivechatRequests, id livechat.LicenseID, localURL string) *app {
	return &app{
		lcHTTP:    lcHTTP,
		licenseID: id,
		agents:    &agents{},
		webhooks:  make(map[string]*webhookDetails),
		localURL:  localURL,
	}
}

func (a *app) RegisterAction(ctx context.Context, actions ...string) error {
	clientID, err := auth.GetClientID(ctx)
	if err != nil {
		return fmt.Errorf("bot: register_action: %w", err)
	}

	for _, action := range actions {
		payload := &web.RegisterWebhookRequest{
			ClientID:  clientID,
			SecretKey: "random secret key",
			URL:       fmt.Sprintf("%s/webhooks/%s", a.localURL, action),
			Action:    action,
			Type:      "license",
		}

		webhookResponse, err := a.lcHTTP.RegisterWebhook(ctx, payload)
		if err != nil {
			return fmt.Errorf("bot: register_action: %w", err)
		}

		log.WithField("action", action).WithField("url", fmt.Sprintf("%s/webhooks/%s", a.localURL, action)).Debug("Webhook registered")

		a.webhooks[action] = &webhookDetails{id: webhookResponse.ID}
	}

	return nil
}

func (a *app) UnregisterActions(ctx context.Context) error {
	clientID, err := auth.GetClientID(ctx)
	if err != nil {
		return fmt.Errorf("bot: register_action: %w", err)
	}

	wg := sync.WaitGroup{}

	for _, agent := range a.agents.agents {
		wg.Add(1)
		go func(agentID livechat.AgentID) {
			defer wg.Done()
			if err := newBotFactory(a.lcHTTP).disableBot(ctx, agentID); err == nil {
				log.WithField("license_id", a.licenseID).Debug("Bots are offline")
			}
		}(agent.ID)
	}

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

func (a *app) TransferChat(ctx context.Context, msg *rtm.PushIncomingChat, data ...RedirectData) error {
	agent, err := a.agents.FindByChatExclude(msg.Payload.Chat.ID)
	if err != nil {
		return fmt.Errorf("bot: transfer_chat action: %w", err)
	}

	if _, err = a.lcHTTP.TransferChat(ctx, buildTransferChatMessage(msg.Payload.Chat.ID, agent.ID)); !isTransferChatErrorOk(err) {
		return fmt.Errorf("bot: transfer_chat action: %w", err)
	}

	agent.chats = append(agent.chats, msg.Payload.Chat.ID)
	return nil
}

func (a *app) IncomingEvent(ctx context.Context, msg *rtm.PushIncomingMessage, data ...RedirectData) error {
	if msg.Payload.Event.Type != "message" {
		return nil
	}

	if len(data) > 0 && data[0].AppAuthorID == msg.Payload.Event.AuthorID {
		return nil
	}

	agent, err := a.agents.FindByChat(msg.Payload.ChatID)
	if err != nil {
		if err := a.TransferChat(ctx, mapIncomingEventIntoTransferChat(msg)); err != nil {
			return fmt.Errorf("bot: incoming_event action: %w", err)
		}
		return a.IncomingEvent(ctx, msg)
	}

	_, err = a.lcHTTP.SendEvent(auth.WithAuthorID(ctx, agent.ID), &web.SendEventRequest{
		ChatID: msg.Payload.ChatID,
		Event: web.Event{
			Text:       messages.Talk(msg.Payload.Event.Text),
			Type:       "message",
			Recipients: "all",
		},
	})

	if err != nil {
		return fmt.Errorf("bot: incoming_event action: %w", err)
	}
	return nil
}

func buildTransferChatMessage(chatID livechat.ChatID, agentID livechat.AgentID) *web.TransferChatRequest {
	return &web.TransferChatRequest{
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

func mapIncomingEventIntoTransferChat(msg *rtm.PushIncomingMessage) *rtm.PushIncomingChat {
	return &rtm.PushIncomingChat{
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
