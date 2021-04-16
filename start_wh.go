package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/livechat/onboarding/bot"
	"github.com/livechat/onboarding/bot/bot_webhooks"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
)

func StartWebhooks(ctx context.Context, cfg *config, config *appMethodConfig) bot.BotManager {
	// LIVECHAT SERVICES
	lcHTTP := web.New(config.httpClient, cfg.URL.HTTP)
	bot := bot_webhooks.New(lcHTTP, cfg.URL.Local)

	config.router.Post("/webhooks/incoming_event", handleIncomingMsg(ctx, bot, cfg, func() rtm.Push {
		return &rtm.PushIncomingMessage{}
	}))
	config.router.Post("/webhooks/incoming_chat", handleIncomingMsg(ctx, bot, cfg, func() rtm.Push {
		return &rtm.PushIncomingChat{}
	}))

	return bot
}

func handleIncomingMsg(ctx context.Context, bot bot_webhooks.Manager, cfg *config, body func() rtm.Push) http.HandlerFunc {
	additionalData := bot_webhooks.RedirectData{
		AppAuthorID: cfg.Credentials.AuthorID,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyMsg := body()
		if err := json.NewDecoder(r.Body).Decode(&bodyMsg); err != nil {
			sendError(w, err)
			return
		}

		if err := bot.Redirect(ctx, bodyMsg, additionalData); err != nil {
			sendError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
