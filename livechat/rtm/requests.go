package rtm

import (
	"context"
	"fmt"

	"github.com/livechat/onboarding/livechat"
	"github.com/sirupsen/logrus"
)

func (c *livechatClient) SendPing(ctx context.Context) error {
	req := payload{
		"action":  "ping",
		"payload": nil,
	}

	return c.WriteJSON(req)
}

func (c *livechatClient) SendLogin(ctx context.Context, opts *LoginRequest) error {
	req := payload{
		"action": "login",
		"payload": map[string]interface{}{
			"token": fmt.Sprintf("Basic %s", opts.Token),
			"away":  true,
		},
	}

	return c.WriteJSON(req)
}

func (c *livechatClient) SendEvent(ctx context.Context, chatID livechat.ChatID, agentID livechat.AgentID, msg string) error {
	req := payload{
		"action": "send_event",
		"payload": map[string]interface{}{
			"chat_id": chatID,
			"event": map[string]string{
				"type":       "message",
				"text":       msg,
				"recipients": "all",
			},
		},
	}

	logrus.WithField("req", req).Debug("Send event payload")

	return c.WriteJSON(req)
}

func (c *livechatClient) SendTransferChat(ctx context.Context, chatID livechat.ChatID, agentID []string) error {
	req := payload{
		"action": "transfer_chat",
		"payload": map[string]interface{}{
			"id": chatID,
			"target": map[string]interface{}{
				"type": "agent",
				"ids":  agentID,
			},
			"force": true,
		},
	}

	return c.WriteJSON(req)
}

func (c *livechatClient) Read(ctx context.Context) (<-chan []byte, <-chan error) {
	msgHandler := make(chan []byte, 1)
	errHandler := make(chan error, 1)

	go func() {
		defer func() {
			close(msgHandler)
			close(errHandler)
		}()

		for {
			_, incomingMsg, err := c.conn.ReadMessage()
			if err != nil {
				errHandler <- fmt.Errorf("rtm_client: %w", err)
				return
			}

			select {
			case <-ctx.Done():
				errHandler <- fmt.Errorf("rtm_client: %w", ctx.Err())
				return
			case msgHandler <- incomingMsg:
			}
		}
	}()

	return msgHandler, errHandler
}
