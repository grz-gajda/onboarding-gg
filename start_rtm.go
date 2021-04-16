package main

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/livechat/onboarding/bot/bot_rtm"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
	log "github.com/sirupsen/logrus"
)

func StartRTM(ctx context.Context, cfg *config, config *appMethodConfig) BotManager {
	// LIVECHAT SERVICES
	lcHTTP := web.New(config.httpClient, cfg.URL.HTTP)
	lcRTM, err := rtm.New(ctx, &websocket.Dialer{HandshakeTimeout: 5 * time.Second}, cfg.URL.WS)
	if err != nil {
		log.WithError(err).Panic("Cannot initialize connection to LiveChat")
	}

	if err := lcRTM.Login(ctx); err != nil {
		log.WithError(err).Panic("Cannot authorize connection to LiveChat")
	}

	go lcRTM.Ping(ctx)

	return bot_rtm.New(lcHTTP, lcRTM)
}
