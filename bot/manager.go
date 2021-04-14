package bot

import (
	"context"
	"fmt"
	"sync"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
	log "github.com/sirupsen/logrus"
)

type manager struct {
	lcHTTP web.LivechatRequests
	lcRTM  rtm.LivechatCommunicator

	apps apps
}

func (m *manager) InstallApp(ctx context.Context, id livechat.LicenseID) error {
	app := newApp(m.lcHTTP, m.lcRTM, id)
	if err := m.apps.Register(app); err != nil {
		return fmt.Errorf("bot: install_app: %w", err)
	}

	if err := app.FetchBots(ctx); err != nil {
		return fmt.Errorf("bot: install_app: %w", err)
	}

	log.WithField("license_id", id).Info("Installed application")
	return nil
}

func (m *manager) UninstallApp(ctx context.Context, id livechat.LicenseID) error {
	uninstalledApp := m.apps.Unregister(id)
	if uninstalledApp == nil {
		return fmt.Errorf("bot: app (license id: %v) is not registered", id)
	}

	for _, agent := range uninstalledApp.agents.agents {
		agent.closeFn()
	}

	log.WithField("license_id", id).Info("Uninstalled application")
	return nil
}

func (m *manager) Destroy(ctx context.Context) {
	wg := &sync.WaitGroup{}

	for _, app := range m.apps.apps {
		wg.Add(1)

		go func(id livechat.LicenseID) {
			defer wg.Done()
			m.UninstallApp(ctx, id)
		}(app.licenseID)
	}

	wg.Wait()
}

func (m *manager) JoinChat(ctx context.Context, license livechat.LicenseID, chat livechat.ChatID) error {
	agents := []string{}
	app, err := m.apps.Find(license)
	if err != nil {
		return err
	}

	for _, a := range app.agents.agents {
		agents = append(agents, string(a.ID))
	}

	return m.lcRTM.SendTransferChat(ctx, chat, agents)
}
