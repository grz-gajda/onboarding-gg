package agents

import (
	"sync"

	"github.com/livechat/onboarding/livechat"
)

type Agent struct {
	ID livechat.AgentID

	mu    sync.Mutex
	chats []livechat.ChatID
}

func (a *Agent) RegisterChat(chatIDs ...livechat.ChatID) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	chatsToRegister := []livechat.ChatID{}
	chatsToRegister = append(chatsToRegister, a.chats...)

out:
	for _, chatID := range chatIDs {
		for _, registeredChatID := range chatsToRegister {
			if registeredChatID == chatID {
				continue out
			}
		}
		chatsToRegister = append(chatsToRegister, chatID)
	}

	a.chats = chatsToRegister
	return nil
}
