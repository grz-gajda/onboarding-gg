package rtm

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/livechat/onboarding/livechat/auth"
	log "github.com/sirupsen/logrus"
)

type payload map[string]interface{}

type livechatClient struct {
	conn *websocket.Conn
}

func (c *livechatClient) Ping(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	errHandler := make(chan error)

	defer func() {
		ticker.Stop()
		close(errHandler)

		c.Close()
	}()

	isPingOK := func() {
		if err := c.SendPing(ctx); err != nil {
			errHandler <- err
			return
		}
		log.Debug("Sent ping to LiveChat")
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("rtm_client: ping action: %w", ctx.Err())
		case err := <-errHandler:
			log.WithError(err).Warn("Agent stopped working")
			return fmt.Errorf("rtm_client: ping action: %w", err)
		case <-ticker.C:
			go isPingOK()
		}
	}
}

func (c *livechatClient) Login(ctx context.Context) error {
	authToken, err := auth.GetAuthToken(ctx)
	if err != nil {
		return fmt.Errorf("rtm_client: login action: %w", err)
	}

	err = c.SendLogin(ctx, &LoginRequest{Token: authToken})
	if err != nil {
		return fmt.Errorf("rtm_client: login action: %w", err)
	}

	log.WithContext(ctx).Debug("Sent request for authorization")

	return nil
}

func (c *livechatClient) Close() error {
	return c.conn.Close()
}

func (c *livechatClient) WriteJSON(v interface{}) error {
	if err := c.conn.WriteJSON(v); err != nil {
		return fmt.Errorf("rtm_client: %w", err)
	}
	return nil
}
