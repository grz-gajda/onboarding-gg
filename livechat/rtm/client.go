package rtm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type payload map[string]interface{}

type livechatClient struct {
	conn *websocket.Conn
}

func (c *livechatClient) Call(ctx context.Context, payload payload) (<-chan []byte, error) {
	msg := make(chan []byte)
	if err := c.conn.WriteJSON(payload); err != nil {
		close(msg)
		return msg, fmt.Errorf("rtm_client: %w", err)
	}

	go c.FindByAction(ctx, "ping", msg)

	return msg, nil
}

func (c *livechatClient) FindByAction(ctx context.Context, actionName string, msg chan<- []byte) error {
	type actionMsg struct {
		Action string `json:"action"`
	}

	for {
		_, incomingMsg, err := c.conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("rtm_client: %w", err)
		}

		action := actionMsg{}
		if err := json.Unmarshal(incomingMsg, &action); err != nil {
			return fmt.Errorf("rtm_client: %w", err)
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("rtm_client: %w", ctx.Err())
		default:
			if actionName == action.Action {
				msg <- incomingMsg
				close(msg)
				return nil
			}
		}
	}
}

func (c *livechatClient) Close() error {
	return c.conn.Close()
}
