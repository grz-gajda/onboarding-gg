package agents

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/web"
)

func Initialize(ctx context.Context, lcHTTP web.LivechatRequests) (Agents, error) {
	agents := NewCollection()
	_, err := createBot(ctx, lcHTTP)
	if err != nil {
		return agents, fmt.Errorf("bot_factory: %w", err)
	}

	bots, err := fetchBots(ctx, lcHTTP)
	if err != nil {
		return agents, fmt.Errorf("bot_factory: %w", err)
	}
	if bots.Len() == 0 {
		return agents, errors.New("bot_factory: received empty list of bots")
	}

	insideBots, unlockBots := bots.Get()
	defer unlockBots()
	for _, insideBot := range insideBots {
		if err := enableBot(ctx, lcHTTP, insideBot.ID); err != nil {
			go removeBot(ctx, lcHTTP, insideBot.ID)
			continue
		}
		agents.Register(insideBot)
	}

	if agents.Len() == 0 {
		return agents, fmt.Errorf("bot_factory: received empty list of bots")
	}

	return agents, nil
}

func Terminate(ctx context.Context, lcHTTP web.LivechatRequests, bots Agents) error {
	agentsInside, unlock := bots.Get()
	defer unlock()

	wg := sync.WaitGroup{}
	wg.Add(bots.Len())

	for _, agent := range agentsInside {
		go func(bot *Agent) {
			defer wg.Done()
			disableBot(ctx, lcHTTP, bot.ID)
		}(agent)
	}

	wg.Wait()
	return nil
}

func createBot(ctx context.Context, lcHTTP web.LivechatRequests) (*Agent, error) {
	response, err := lcHTTP.CreateBot(ctx, &livechat.CreateBotRequest{Name: "OnboardingGG (bot created by app)"})
	if err != nil {
		return nil, fmt.Errorf("create_bot: %w", err)
	}

	return NewAgent(response.ID), nil
}

func fetchBots(ctx context.Context, lcHTTP web.LivechatRequests) (Agents, error) {
	collection := NewCollection()
	botsResponse, err := lcHTTP.ListBots(ctx, &livechat.ListBotsRequest{All: true})
	if err != nil {
		return collection, fmt.Errorf("fetch_bots: %w", err)
	}
	if len(botsResponse) == 0 {
		return collection, fmt.Errorf("fetch_bots: empty list")
	}

	for _, botID := range botsResponse {
		collection.Register(NewAgent(botID.ID))
	}

	return collection, nil
}

func removeBot(ctx context.Context, lcHTTP web.LivechatRequests, botID livechat.AgentID) error {
	_, err := lcHTTP.DeleteBot(ctx, &livechat.DeleteBotRequest{ID: botID})
	if err != nil {
		return fmt.Errorf("remove_bot: %w", err)
	}
	return nil
}

func enableBot(ctx context.Context, lcHTTP web.LivechatRequests, botID livechat.AgentID) error {
	_, err := lcHTTP.SetRoutingStatus(ctx, &livechat.SetRoutingStatusRequest{
		Status:  "accepting_chats",
		AgentID: botID,
	})

	if err != nil {
		return fmt.Errorf("enable_bot: %w", err)
	}

	return nil
}

func disableBot(ctx context.Context, lcHTTP web.LivechatRequests, botID livechat.AgentID) error {
	_, err := lcHTTP.SetRoutingStatus(ctx, &livechat.SetRoutingStatusRequest{
		Status:  "offline",
		AgentID: botID,
	})

	if err != nil {
		return fmt.Errorf("disable_bot: %w", err)
	}
	return nil
}
