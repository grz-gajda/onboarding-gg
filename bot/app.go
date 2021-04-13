package bot

import (
	"context"
	"fmt"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
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
		return fmt.Errorf("app_manager create_bot: %w", err)
	}

	response, err := a.lcHTTP.CreateBot(ctx, &web.CreateBotRequest{
		Name:     "OnboardingGG",
		ClientID: clientID,
	})
	if err != nil {
		return fmt.Errorf("app_manager create_bot: %w", err)
	}

	agent := newAgent(response.ID, a.lcRTM)
	a.agents.Register(agent)

	return nil
}

func (a *app) FetchBots(ctx context.Context) error {
	botsResponse, err := a.lcHTTP.ListBots(ctx, &web.ListBotsRequest{All: true})
	if err != nil {
		return fmt.Errorf("app_manager fetch_bots: %w", err)
	}
	if len(botsResponse) == 0 {
		return a.CreateBot(ctx)
	}

	for _, botID := range botsResponse {
		agent := newAgent(botID.ID, a.lcRTM)
		a.agents.Register(agent)
	}

	return nil
}
