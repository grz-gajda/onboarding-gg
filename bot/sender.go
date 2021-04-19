package bot

import (
	"context"

	"github.com/livechat/onboarding/livechat"
)

type sender struct {
	client      SenderClient
	appAuthorID string
}

func NewSender(client SenderClient, authorID string) Sender {
	return &sender{client: client, appAuthorID: authorID}
}

func (s *sender) Talk(ctx context.Context, chatID livechat.ChatID, msg *livechat.PushIncomingMessage) error {
	if msg.Payload.Event.Type != "message" {
		return nil
	}
	if msg.Payload.Event.AuthorID == s.appAuthorID {
		return nil
	}

	return s.client.SendEvent(ctx, &livechat.Event{
		ChatID: chatID,
		Event: livechat.EventMessage{
			Type: livechat.EventTypeMessage,
			Text: "Hello world!",
		},
	})
}
