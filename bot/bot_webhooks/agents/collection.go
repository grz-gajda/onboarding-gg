package agents

import (
	"errors"
	"fmt"
	"sync"

	"github.com/livechat/onboarding/livechat"
)

type collection struct {
	agents []*Agent
	mu     sync.Mutex
}

func (a *collection) Register(bot *Agent) error {
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

func (a *collection) Unregister(agentID livechat.AgentID) (*Agent, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	newAgents := []*Agent{}
	var unregisteredAgent *Agent

	for _, agent := range a.agents {
		if agent.ID != agentID {
			newAgents = append(newAgents, agent)
		} else {
			unregisteredAgent = agent
		}
	}

	if nil == unregisteredAgent {
		return nil, errors.New("bot: agent cannot be found")
	}

	a.agents = newAgents
	return unregisteredAgent, nil
}

func (a *collection) Len() int { return len(a.agents) }

func (a *collection) Get() ([]*Agent, AgentsUnlock) {
	a.mu.Lock()

	return a.agents, a.mu.Unlock
}

func (a *collection) FindByChat(chatID livechat.ChatID) (*Agent, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, agent := range a.agents {
		isChatArchived := func(c livechat.ChatID) bool {
			for _, chat := range agent.removed {
				if chat == c {
					return true
				}
			}
			return false
		}

		for _, chat := range agent.chats {
			if chat == chatID && !isChatArchived(chat) {
				return agent, nil
			}
		}
	}

	return nil, fmt.Errorf("bot: agent cannot be found")
}

func (a *collection) FindByChatExclude(chatID livechat.ChatID) (*Agent, error) {
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

func (a *collection) FindByID(agentID livechat.AgentID) (*Agent, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, agent := range a.agents {
		if agent.ID == agentID {
			return agent, nil
		}
	}

	return nil, fmt.Errorf("bot: agent cannot be found")
}
