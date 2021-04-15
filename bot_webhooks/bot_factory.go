package bot_webhooks

import (
	"context"
	"errors"
	"fmt"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/web"
	"github.com/sirupsen/logrus"
)

type botFactory struct {
	lcHTTP web.LivechatRequests
}

func newBotFactory(lcHTTP web.LivechatRequests) *botFactory {
	return &botFactory{lcHTTP: lcHTTP}
}

func (f *botFactory) Initialize(ctx context.Context) ([]*agent, error) {
	_, err := f.createBot(ctx)
	if err != nil {
		return []*agent{}, fmt.Errorf("bot_factory: %w", err)
	}

	bots, err := f.fetchBots(ctx)
	if err != nil {
		return []*agent{}, fmt.Errorf("bot_factory: %w", err)
	}
	if len(bots) == 0 {
		return []*agent{}, errors.New("bot_factory: received empty list of bots")
	}

	var enabledBots []*agent
	for _, bot := range bots {
		if err := f.enableBot(ctx, bot.ID); err != nil {
			f.removeBot(ctx, bot.ID)
			continue
		}

		enabledBots = append(enabledBots, bot)
	}

	if len(enabledBots) == 0 {
		return []*agent{}, fmt.Errorf("bot_factory: received empty list of bots")
	}

	logrus.WithField("registered_bots", len(enabledBots)).Debug("Bots have been registered and enabled")

	return enabledBots, nil
}

func (f *botFactory) createBot(ctx context.Context) (*agent, error) {
	clientID, err := auth.GetClientID(ctx)
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	response, err := f.lcHTTP.CreateBot(ctx, &web.CreateBotRequest{
		Name:     "OnboardingGG",
		ClientID: clientID,
	})
	if err != nil {
		return nil, fmt.Errorf("create_bot: %w", err)
	}

	return newAgent(response.ID, f.lcHTTP), nil
}

func (f *botFactory) fetchBots(ctx context.Context) ([]*agent, error) {
	botsResponse, err := f.lcHTTP.ListBots(ctx, &web.ListBotsRequest{All: true})
	if err != nil {
		return []*agent{}, fmt.Errorf("fetch_bots: %w", err)
	}
	if len(botsResponse) == 0 {
		return []*agent{}, fmt.Errorf("fetch_bots: empty list")
	}

	var bots []*agent
	for _, botID := range botsResponse {
		bots = append(bots, newAgent(botID.ID, f.lcHTTP))
	}

	return bots, nil
}

func (f *botFactory) removeBot(ctx context.Context, botID livechat.AgentID) error {
	_, err := f.lcHTTP.DeleteBot(ctx, &web.DeleteBotRequest{ID: botID})
	if err != nil {
		return fmt.Errorf("remove_bot: %w", err)
	}
	return nil
}

func (f *botFactory) enableBot(ctx context.Context, botID livechat.AgentID) error {
	_, err := f.lcHTTP.SetRoutingStatus(ctx, &web.SetRoutingStatusRequest{
		Status:  "accepting_chats",
		AgentID: botID,
	})

	if err != nil {
		return fmt.Errorf("enable_bot: %w", err)
	}

	return nil
}

func (f *botFactory) disableBot(ctx context.Context, botID livechat.AgentID) error {
	_, err := f.lcHTTP.SetRoutingStatus(ctx, &web.SetRoutingStatusRequest{
		Status:  "offline",
		AgentID: botID,
	})

	if err != nil {
		return fmt.Errorf("disable_bot: %w", err)
	}
	return nil
}
