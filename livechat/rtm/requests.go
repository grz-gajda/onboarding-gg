package rtm

import (
	"context"
	"fmt"

	"github.com/livechat/onboarding/livechat"
)

func (c *livechatClient) SendPing(ctx context.Context) error {
	req := payload{
		"action":  "ping",
		"payload": nil,
	}

	return c.conn.WriteJSON(req)
}

func (c *livechatClient) SendLogin(ctx context.Context, opts *LoginRequest) error {
	req := payload{
		"action": "login",
		"payload": map[string]interface{}{
			"token": fmt.Sprintf("Basic %s", opts.Token),
			"away":  true,
		},
	}

	return c.conn.WriteJSON(req)
}

func (c *livechatClient) SendTransferChat(ctx context.Context, chatID livechat.ChatID, agents []livechat.AgentID) error {
	req := payload{
		"action": "transfer_chat",
		"payload": map[string]interface{}{
			"chat_id": chatID,
			"target": map[string]interface{}{
				"type": "agent",
				"ids":  agents,
			},
		},
	}

	return c.conn.WriteJSON(req)
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
				errHandler <- err
				return
			}

			select {
			case <-ctx.Done():
				errHandler <- ctx.Err()
				return
			case msgHandler <- incomingMsg:
			}
		}
	}()

	return msgHandler, errHandler
}
