package bot_webhooks

import (
	"context"

	"github.com/livechat/onboarding/bot"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/web"
)

var WebhookEvents = []string{
	"incoming_chat",
	"incoming_event",
}

type Manager interface {
	bot.BotManager
	Redirect(context.Context, livechat.Push) error
}

func New(lcHTTP web.LivechatRequests, localURL string, authorID string) Manager {
	return &manager{
		lcHTTP:   lcHTTP,
		localURL: localURL,
		apps:     &apps{},
		sender:   bot.NewSender(&senderClient{lcHTTP: lcHTTP}, authorID),
	}
}

type senderClient struct {
	lcHTTP web.LivechatRequests
}

func (s *senderClient) SendEvent(ctx context.Context, event *livechat.Event) error {
	_, err := s.lcHTTP.SendEvent(ctx, event)
	return err
}
