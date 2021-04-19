package agents

import (
	"sync"

	"github.com/livechat/onboarding/livechat"
)

type AgentsUnlock func()

type Agents interface {
	Register(*Agent) error
	Unregister(livechat.AgentID) (*Agent, error)

	Len() int
	Get() ([]*Agent, AgentsUnlock)

	FindByChat(livechat.ChatID) (*Agent, error)
	FindByChatExclude(livechat.ChatID) (*Agent, error)
}

func NewCollection() Agents {
	return &collection{
		agents: []*Agent{},
		mu:     sync.Mutex{},
	}
}

func NewAgent(id livechat.AgentID) *Agent {
	return &Agent{
		ID:    id,
		chats: []livechat.ChatID{},
		mu:    sync.Mutex{},
	}
}
