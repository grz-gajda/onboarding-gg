package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/sirupsen/logrus"
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

func (a *agent) Start(ctx context.Context) error {
	msgHandler, errHandler := a.Conn.Read(ctx)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("bot agent: start action: %w", ctx.Err())
		case msg := <-msgHandler:
			var body interface{}
			if err := json.Unmarshal(msg, &body); err == nil {
				logrus.WithField("message", body).Debug("Received message")
			}
		case err := <-errHandler:
			return fmt.Errorf("bot agent: start action: %w", err)
		case <-a.closeCh:
			return errors.New("bot agent: start action: agent terminated")
		}
	}
}
