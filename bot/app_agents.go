package bot

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
			return fmt.Errorf("agent (id: %v) is already registered", bot.ID)
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

func (a *agents) FindBy(id livechat.AgentID) (*agent, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, agent := range a.agents {
		if agent.ID == id {
			return agent, nil
		}
	}
	return nil, errors.New("agents: agent cannot be found")
}
