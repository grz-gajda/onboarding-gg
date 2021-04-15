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

	url := "http://localhost:8081"
	if opts.CustomURL != "" {
		url = opts.CustomURL
	}

	payload := &web.RegisterWebhookRequest{
		ClientID:  clientID,
		SecretKey: "random secret key",
		URL:       fmt.Sprintf("%s/webhooks/%s", url, action),
		Action:    action,
		Type:      "license",
	}

	webhookResponse, err := a.lcHTTP.RegisterWebhook(ctx, payload)
	if err != nil {
		return fmt.Errorf("bot: register_action: %w", err)
	}

	log.WithField("action", action).WithField("url", fmt.Sprintf("%s/webhooks/%s", url, action)).Debug("Webhook registered")

	a.webhooks[action] = &webhookDetails{id: webhookResponse.ID}
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
			newBotFactory(a.lcHTTP).disableBot(ctx, agentID)
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
		Force: true,
	})

	if err != nil {
		if !strings.Contains(err.Error(), "One or more of requested agents are already present in the chat") {
			return fmt.Errorf("bot: transfer_chat action: %w", err)
		}
	}

	agent.chats = append(agent.chats, msg.Payload.Chat.ID)
	return nil
}

func (a *app) IncomingEvent(ctx context.Context, msg *rtm.PushIncomingMessage, data ...RedirectData) error {
	log.WithField("agents", a.agents).WithField("chat_id", msg.Payload.ChatID).Debug("IncomingEvent action")
	if msg.Payload.Event.Type != "message" {
		return nil
	}

	log.WithField("data", data).WithField("author_id", msg.Payload.Event.AuthorID).Debug("Comparing 'author_id'`s of messages")
	if len(data) > 0 && data[0].AppAuthorID == msg.Payload.Event.AuthorID {
		return nil
	}

	agent, err := a.agents.FindByChat(msg.Payload.ChatID)
	if err != nil {
		err := a.TransferChat(ctx, &rtm.PushIncomingChat{
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
		})

		if err != nil {
			return fmt.Errorf("bot: incoming_event action: %w", err)
		}

		return a.IncomingEvent(ctx, msg)
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
