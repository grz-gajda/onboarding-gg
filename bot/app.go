package bot

import (
	"context"
	"fmt"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
	log "github.com/sirupsen/logrus"
)

type app struct {
	lcHTTP    web.LivechatRequests
	lcRTM     rtm.LivechatCommunicator
	agents    *agents
	licenseID livechat.LicenseID
}

func newApp(lcClient web.LivechatRequests, rtmClient rtm.LivechatCommunicator, licenseID livechat.LicenseID) *app {
	return &app{
		agents:    &agents{},
		lcHTTP:    lcClient,
		lcRTM:     rtmClient,
		licenseID: licenseID,
	}
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

	agent := newAgent(response.ID, a.lcRTM)
	if err := a.agents.Register(agent); err != nil {
		return fmt.Errorf("bot: create_bot: %w", err)
	}

	go agent.Start(ctx)

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
		agent := newAgent(botID.ID, a.lcRTM)
		if err := a.agents.Register(agent); err != nil {
			return fmt.Errorf("bot: fetch_bots: %w", err)
		}

		go agent.Start(ctx)
	}

	log.WithField("license_id", a.licenseID).WithField("amount", len(botsResponse)).Infof("Bot has been registered")
	return nil
}
