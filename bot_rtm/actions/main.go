package actions

import (
	"context"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/rtm"
)

type Actions interface {
	IncomingEvent(context.Context, livechat.AgentID, []byte) error
}

func New(lcRTM rtm.LivechatCommunicator) Actions {
	return &actions{lcRTM: lcRTM}
}

type actions struct {
	lcRTM rtm.LivechatCommunicator
}
