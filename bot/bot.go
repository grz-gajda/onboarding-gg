package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/livechat/onboarding/bot/actions"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/rtm"
	log "github.com/sirupsen/logrus"
)

type agent struct {
	ID   livechat.AgentID
	Conn rtm.LivechatCommunicator

	closeCh chan struct{}
	closeFn func()
}

func newAgent(ID livechat.AgentID, conn rtm.LivechatCommunicator) *agent {
	closeCh := make(chan struct{}, 1)
	closeFn := func() {
		closeCh <- struct{}{}
		close(closeCh)
	}

	return &agent{
		ID:      ID,
		Conn:    conn,
		closeCh: closeCh,
		closeFn: closeFn,
	}
}

func (a *agent) Start(ctx context.Context) (stopCause error) {
	msgHandler, errHandler := a.Conn.Read(ctx)
	defer func() {
		log.WithField("agent_id", a.ID).WithError(stopCause).Debug("Stopped agent")
	}()

	for {
		select {
		case <-ctx.Done():
			stopCause = fmt.Errorf("bot: start action: %w", ctx.Err())
		case msg := <-msgHandler:
			if err := a.HandleMsg(ctx, msg); err != nil {
				stopCause = err
			}
		case err := <-errHandler:
			stopCause = fmt.Errorf("bot: start action: %w", err)
		case <-a.closeCh:
			stopCause = errors.New("bot: start action: agent terminated")
		}

		if stopCause != nil {
			return
		}
	}
}

func (a *agent) HandleMsg(ctx context.Context, msg []byte) error {
	var body map[string]interface{}
	if err := json.Unmarshal(msg, &body); err != nil {
		return fmt.Errorf("bot: cannot unmarshal incoming msg: %w", err)
	}

	action, ok := body["action"]
	if !ok {
		return errors.New("bot: incoming msg does not contain 'action' field")
	}

	switch action {
	case "incoming_event":
		return actions.New(a.Conn).IncomingEvent(ctx, a.ID, msg)
	default:
		log.WithField("body", body).Debug("Incoming unknown message")
	}

	return nil
}
