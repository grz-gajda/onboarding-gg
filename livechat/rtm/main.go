package rtm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/livechat/onboarding/livechat"
	log "github.com/sirupsen/logrus"
)

type Client interface {
	DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)
}

type LivechatCommunicator interface {
	// STATE
	Close() error

	// MANAGE
	Ping(context.Context) error
	Login(context.Context) error

	// ACTIONS
	SendEvent(context.Context, livechat.ChatID, livechat.AgentID, string) error
	SendTransferChat(context.Context, livechat.ChatID, []string) error

	// READERS
	Read(context.Context) (<-chan []byte, <-chan error)
}

func New(ctx context.Context, client Client, url string) (LivechatCommunicator, error) {
	conn, _, err := client.DialContext(ctx, url, nil)
	if err != nil {
		return &livechatClient{}, fmt.Errorf("rtm_client: %w", err)
	}

	log.WithContext(ctx).Info("Initialized connection with LiveChat RTM")

	return &livechatClient{conn: conn}, nil
}

type SendEvent struct {
	Type       string `json:"type"`
	Text       string `json:"text"`
	Recipients string `json:"recipients"`
}

type Action struct {
	Action string `json:"action"`
}

type LoginRequest struct {
	Token string
}
