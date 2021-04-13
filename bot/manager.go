package bot

import (
	"context"
	"errors"
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
		return err
	}

	if err := app.FetchBots(ctx); err != nil {
		return err
	}

	log.WithField("license_id", id).Info("Installed application")
	return nil
}

func (m *manager) UninstallApp(ctx context.Context, id livechat.LicenseID) error {
	uninstalledApp := m.apps.Unregister(id)
	if uninstalledApp == nil {
		return fmt.Errorf("app (license id: %v) is not registered", id)
	}

	for _, agent := range uninstalledApp.agents.agents {
		agent.closeFn()
		if _, err := m.lcHTTP.DeleteBot(ctx, &web.DeleteBotRequest{ID: agent.ID}); err != nil {
			log.WithField("agent_id", agent.ID).WithContext(ctx).WithError(err).Error("Cannot remove agent from LiveChat")
		}
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
	return errors.New("not implemented yet")
	// m.mu.Lock()
	// var app *app
	// for _, a := range m.apps {
	// 	if a.licenseID == license {
	// 		app = a
	// 		break
	// 	}
	// }
	// m.mu.Unlock()

	// agents := []livechat.AgentID{}
	// for _, agent := range app.bots {
	// 	agents = append(agents, agent.ID)
	// }

	// return nil
}
