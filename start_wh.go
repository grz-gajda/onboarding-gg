package main

import (
	"encoding/json"
	"net/http"

	"github.com/livechat/onboarding/bot"
	"github.com/livechat/onboarding/bot/bot_webhooks"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/web"
)

func StartWebhooks(cfg *config, config *appMethodConfig) bot.BotManager {
	// LIVECHAT SERVICES
	lcHTTP := web.New(config.httpClient, cfg.URL.HTTP)
	bot := bot_webhooks.New(lcHTTP, cfg.URL.Local, cfg.Credentials.AuthorID)

	config.router.Post("/webhooks/incoming_event", handleIncomingMsg(bot, cfg, func() livechat.Push {
		return &livechat.PushIncomingMessage{}
	}))
	config.router.Post("/webhooks/incoming_chat", handleIncomingMsg(bot, cfg, func() livechat.Push {
		return &livechat.PushIncomingChat{}
	}))
	config.router.Post("/webhooks/user_added_to_chat", handleIncomingMsg(bot, cfg, func() livechat.Push {
		return &livechat.PushUserAddedToChat{}
	}))

	return bot
}

func handleIncomingMsg(bot bot_webhooks.Manager, cfg *config, body func() livechat.Push) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := auth.WithClientID(r.Context(), cfg.Credentials.ClientID)

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
