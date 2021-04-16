package bot_webhooks

import (
	"errors"
	"fmt"
	"sync"

	"github.com/livechat/onboarding/livechat"
)

type agents struct {
	agents []*agent
	mu     sync.Mutex
}

func (a *agents) Register(bot *agent) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, agent := range a.agents {
		if agent.ID == bot.ID {
			return fmt.Errorf("bot: agent (id: %v) is already registered", bot.ID)
		}
	}

	a.agents = append(a.agents, bot)
	return nil
}

func (a *agents) Unregister(agentID livechat.AgentID) *agent {
	a.mu.Lock()
	defer a.mu.Unlock()

	newAgents := []*agent{}
	var unregisteredAgent *agent

	for _, agent := range a.agents {
		if agent.ID != agentID {
			newAgents = append(newAgents, agent)
		} else {
			unregisteredAgent = agent
		}
	}

	a.agents = newAgents
	return unregisteredAgent
}

func (a *agents) FindByChat(chatID livechat.ChatID) (*agent, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, agent := range a.agents {
		for _, chat := range agent.chats {
			if chat == chatID {
				return agent, nil
			}
		}
	}

	return nil, fmt.Errorf("bot: agent cannot be found")
}

func (a *agents) FindByChatExclude(chatID livechat.ChatID) (*agent, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, agent := range a.agents {
		if len(agent.chats) == 0 {
			return agent, nil
		}
		hasChat := false

		for _, chat := range agent.chats {
			if chat == chatID {
				hasChat = true
			}
		}

		if !hasChat {
			return agent, nil
		}
	}

	return nil, fmt.Errorf("bot: agent cannot be found")
}

func (a *agents) FindBy(id livechat.AgentID) (*agent, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, agent := range a.agents {
		if agent.ID == id {
			return agent, nil
		}
	}
	return nil, errors.New("bot: agent cannot be found")
}
