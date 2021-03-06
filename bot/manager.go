package bot

import (
	"context"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
)

type BotManager interface {
	Authorize(context.Context, livechat.Client, *auth.AuthorizeCredentials) error
	// InstallApp allows to register new license into memory
	// and fetch existing (or create a new one) bots (agents).
	InstallApp(context.Context, livechat.LicenseID) error
	// UninstallApp removes the whole memory footprint and cancel
	// existing connections to LiveChat.
	UninstallApp(context.Context, livechat.LicenseID) error
	// Destroy does everything what UninstallApp but for every license.
	Destroy(context.Context)
}

type Sender interface {
	Talk(context.Context, livechat.ChatID, *livechat.PushIncomingMessage) error
}
