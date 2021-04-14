package bot_webhooks

import (
	"context"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
)

type Manager interface {
	// InstallApp allows to register new license into memory
	// and fetch existing (or create a new one) bots (agents).
	InstallApp(context.Context, livechat.LicenseID) error
	// UninstallApp removes the whole memory footprint and cancel
	// existing connections to LiveChat.
	UninstallApp(context.Context, livechat.LicenseID) error
	// Destroy does everything what UninstallApp but for every license.
	Destroy(context.Context)

	Redirect(context.Context, rtm.Push) error
}

func New(lcHTTP web.LivechatRequests) Manager {
	return &manager{
		lcHTTP: lcHTTP,
		apps:   &apps{},
	}
}
