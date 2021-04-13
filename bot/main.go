package bot

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/web"
)

type UninstallBot func(context.Context, string) error

type BotManager interface {
	InstallApp(context.Context, livechat.LicenseID) error
	UninstallApp(context.Context, livechat.LicenseID) error
	Destroy(context.Context)
	JoinChat(context.Context, livechat.LicenseID, livechat.ChatID) error
}

func New(websocketURL string, httpClient web.LivechatRequests) BotManager {
	return &manager{
		wsURL:  websocketURL,
		lcHTTP: httpClient,
		dialer: &websocket.Dialer{HandshakeTimeout: 5 * time.Second},
	}
}
