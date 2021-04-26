package agents

import (
	"sync"

	"github.com/livechat/onboarding/livechat"
)

type Agent struct {
	ID livechat.AgentID

	mu      sync.Mutex
	chats   []livechat.ChatID
	removed []livechat.ChatID
}

func (a *Agent) RegisterChat(chatIDs ...livechat.ChatID) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	newRemovedIDs := []livechat.ChatID{}
	chatsToRegister := []livechat.ChatID{}
	chatsToRegister = append(chatsToRegister, a.chats...)

out1:
	for _, chatID := range chatIDs {
		for _, registeredChatID := range chatsToRegister {
			if registeredChatID == chatID {
				continue out1
			}
		}
		chatsToRegister = append(chatsToRegister, chatID)
	}

out2:
	for _, removedID := range a.removed {
		for _, id := range chatIDs {
			if id == removedID {
				continue out2
			}
		}
		newRemovedIDs = append(newRemovedIDs, removedID)
	}

	a.chats = chatsToRegister
	a.removed = newRemovedIDs
	return nil
}

func (a *Agent) UnregisterChat(chatID livechat.ChatID) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	chatsToRegister := []livechat.ChatID{}
	for _, chat := range a.chats {
		if chat != chatID {
			chatsToRegister = append(chatsToRegister, chat)
		} else {
			a.removed = append(a.removed, chat)
		}
	}

	a.chats = chatsToRegister
	return nil
}
