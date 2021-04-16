package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/livechat/onboarding/bot"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/rtm"
	log "github.com/sirupsen/logrus"
)

func (a *actions) IncomingEvent(ctx context.Context, agentID livechat.AgentID, msg []byte) error {
	var payload *rtm.PushIncomingMessage
	if err := json.Unmarshal(msg, &payload); err != nil {
		return fmt.Errorf("actions: incoming_event: %w", err)
	}

	log.WithField("event_author_id", payload.Payload.Event.AuthorID).WithField("agent_id", agentID).Debug("The same author, skipping")

	if payload.Payload.Event.AuthorID == string(agentID) {
		return nil
	}

	return a.lcRTM.SendEvent(ctx, payload.Payload.ChatID, agentID, bot.Talk(payload.Payload.Event.Text))
}
