package bot_webhooks

import (
	"sync"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/web"
)

type agent struct {
	ID   livechat.AgentID
	Conn web.LivechatRequests

	mu    sync.Mutex
	chats []livechat.ChatID
}

func newAgent(id livechat.AgentID, conn web.LivechatRequests) *agent {
	return &agent{ID: id, Conn: conn}
}

func (a *agent) AppendChat(chatID livechat.ChatID) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.chats = append(a.chats, chatID)
}

func (a *agent) RemoveChat(chatID livechat.ChatID) {
	a.mu.Lock()
	defer a.mu.Unlock()

	newChats := []livechat.ChatID{}
	for _, chat := range a.chats {
		if chat != chatID {
			newChats = append(newChats, chat)
		}
	}
	a.chats = newChats
}
