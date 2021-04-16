package bot

import (
	"context"

	"github.com/livechat/onboarding/livechat"
)

type BotManager interface {
	// InstallApp allows to register new license into memory
	// and fetch existing (or create a new one) bots (agents).
	InstallApp(context.Context, livechat.LicenseID) error
	// UninstallApp removes the whole memory footprint and cancel
	// existing connections to LiveChat.
	UninstallApp(context.Context, livechat.LicenseID) error
	// Destroy does everything what UninstallApp but for every license.
	Destroy(context.Context)
}
