package bot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
	"github.com/sirupsen/logrus"
)

type app struct {
	mu        sync.Mutex
	lcHTTP    web.LivechatRequests
	lcRTM     rtm.LivechatCommunicator
	bots      []*agent
	licenseID livechat.LicenseID
}

func newApp(lcClient web.LivechatRequests, rtmClient rtm.LivechatCommunicator, licenseID livechat.LicenseID) *app {
	return &app{
		bots:      []*agent{},
		lcHTTP:    lcClient,
		lcRTM:     rtmClient,
		licenseID: licenseID,
	}
}

func (a *app) Ping(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	errHandler := make(chan error)

	defer func() {
		ticker.Stop()
		close(errHandler)

		logrus.WithField("agent_id", a.licenseID).Debug("Agent stopped")
		a.lcRTM.Close()
	}()

	isPingOK := func() {
		if err := a.lcRTM.SendPing(ctx); err != nil {
			errHandler <- err
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("bot agent: ping action: %w", ctx.Err())
		case err := <-errHandler:
			logrus.WithError(err).Warn("Agent stopped working")
			return fmt.Errorf("bot agent: ping action: %w", err)
		case <-ticker.C:
			go isPingOK()
		}
	}
}

func (a *app) Login(ctx context.Context) error {
	authToken, err := auth.GetAuthToken(ctx)
	if err != nil {
		return fmt.Errorf("bot agent: login action: %w", err)
	}

	err = a.lcRTM.SendLogin(ctx, &rtm.LoginRequest{Token: authToken})
	if err != nil {
		logrus.WithError(err).WithField("token", authToken).Error("Cannot authorize")
		return fmt.Errorf("bot agent: login action: %w", err)
	}

	return nil
}

func (a *app) CreateBot(ctx context.Context, dial *websocket.Dialer, wsURL string) error {
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
	if err := a.registerAgent(ctx, agent); err != nil {
		logrus.WithError(err).WithField("bot id", agent.ID).Error("Cannot register agent into application")
		return fmt.Errorf("app_manager create_bot: %w", err)
	}

	return nil
}

func (a *app) FetchBots(ctx context.Context, dial *websocket.Dialer, wsURL string) error {
	botsResponse, err := a.lcHTTP.ListBots(ctx, &web.ListBotsRequest{All: true})
	if err != nil {
		return fmt.Errorf("app_manager fetch_bots: %w", err)
	}
	if len(botsResponse) == 0 {
		return a.CreateBot(ctx, dial, wsURL)
	}

	for _, botID := range botsResponse {
		agent := newAgent(botID.ID, a.lcRTM)
		if err := a.registerAgent(ctx, agent); err != nil {
			logrus.WithError(err).WithField("bot id", agent.ID).Error("Cannot register agent into application")
			continue
		}
	}

	return nil
}

func (a *app) registerAgent(ctx context.Context, bot *agent) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.bots = append(a.bots, bot)

	return nil
}
