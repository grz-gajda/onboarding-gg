package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/livechat/onboarding/bot_webhooks"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
)

func StartWebhooks(ctx context.Context, cfg *config, config *appMethodConfig) BotManager {
	// LIVECHAT SERVICES
	lcHTTP := web.New(config.httpClient, cfg.URL.HTTP)
	bot := bot_webhooks.New(lcHTTP)

	config.router.Post("/webhooks/incoming_event", handleIncomingMsg(ctx, bot, func() rtm.Push {
		return &rtm.PushIncomingMessage{}
	}))
	config.router.Post("/webhooks/incoming_chat", handleIncomingMsg(ctx, bot, func() rtm.Push {
		return &rtm.PushIncomingChat{}
	}))

	return bot
}

func handleIncomingMsg(ctx context.Context, bot bot_webhooks.Manager, body func() rtm.Push) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyMsg := body()
		if err := json.NewDecoder(r.Body).Decode(&bodyMsg); err != nil {
			sendError(w, err)
			return
		}

		if err := bot.Redirect(ctx, bodyMsg); err != nil {
			sendError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
