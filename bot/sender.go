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

	switch msg.Payload.Event.Text {
	case "Hello":
		return s.client.SendEvent(ctx, livechat.BuildMessage(chatID, "World!"))
	case "Wróc do człowieka":
		return nil
	default:
		return s.client.SendEvent(ctx, livechat.BuildButtonMessage(chatID, "Czy chcesz wrócić do człowieka?", "", "Wróć do człowieka"))
	}
}
