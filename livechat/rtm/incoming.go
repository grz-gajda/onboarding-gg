package rtm

import "github.com/livechat/onboarding/livechat"

type PushIncomingMessage struct {
	Action  string `json:"action"`
	Payload struct {
		ChatID   livechat.ChatID `json:"chat_id"`
		ThreadID string          `json:"thread_id,omitempty"`
		Event    struct {
			Type     string `json:"type"`
			Text     string `json:"text"`
			AuthorID string `json:"author_id"`
		} `json:"event"`
	} `json:"payload"`
}
