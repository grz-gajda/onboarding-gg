package bot_webhooks

import (
	"context"

	"github.com/livechat/onboarding/bot"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
)

type Manager interface {
	bot.BotManager
	Redirect(context.Context, rtm.Push, ...RedirectData) error
}

type RedirectData struct {
	AppAuthorID string
}

func New(lcHTTP web.LivechatRequests, localURL string) Manager {
	return &manager{
		lcHTTP:     lcHTTP,
		localURL:   localURL,
		apps:       &apps{},
		botFactory: newBotFactory(lcHTTP),
	}
}
