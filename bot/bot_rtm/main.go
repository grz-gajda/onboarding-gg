package bot_rtm

import (
	"github.com/livechat/onboarding/bot"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
)

type Manager interface {
	bot.BotManager
}

func New(httpClient web.LivechatRequests, rtmClient rtm.LivechatCommunicator) Manager {
	return &manager{
		lcHTTP: httpClient,
		lcRTM:  rtmClient,
		apps:   apps{},
	}
}
