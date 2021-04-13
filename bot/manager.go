package bot

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
	"github.com/sirupsen/logrus"
)

type manager struct {
	mu     sync.Mutex
	lcHTTP web.LivechatRequests
	dialer *websocket.Dialer
	wsURL  string

	apps []*app
}

func (m *manager) InstallApp(ctx context.Context, id livechat.LicenseID) error {
	m.mu.Lock()

	for _, app := range m.apps {
		if app.licenseID == id {
			m.mu.Unlock()
			return fmt.Errorf("manager: install action: app is already registered")
		}
	}

	conn, err := rtm.New(ctx, m.dialer, m.wsURL)
	if err != nil {
		m.mu.Unlock()
		return fmt.Errorf("manager: install action: %w", err)
	}

	app := newApp(m.lcHTTP, conn, id)
	m.apps = append(m.apps, app)

	m.mu.Unlock()

	go app.Login(ctx)
	go app.Ping(ctx)

	if err := app.FetchBots(ctx, m.dialer, m.wsURL); err != nil {
		logrus.WithError(err).WithContext(ctx).WithField("license id", id).Error("Cannot fetch/create bots for this licesnse")
		return err
	}

	return nil
}

func (m *manager) UninstallApp(ctx context.Context, id livechat.LicenseID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	wg := &sync.WaitGroup{}

	newApps := []*app{}

	for _, app := range m.apps {
		if app.licenseID != id {
			newApps = append(newApps, app)
			continue
		}

		for _, agent := range app.bots {
			wg.Add(1)
			logrus.WithField("bot_id", agent.ID).Debug("Agent is going to be destroyed")

			go func(connClose io.Closer, agentID livechat.AgentID) {
				defer wg.Done()
				defer func() {
					logrus.WithField("bot_id", agentID).Debug("Agent removed")
				}()

				if err := connClose.Close(); err != nil {
					logrus.WithError(err).WithField("bot_id", agentID).Error("Cannot close the connection with LiveChat")
				}

				if _, err := m.lcHTTP.DeleteBot(ctx, &web.DeleteBotRequest{ID: agentID}); err != nil {
					logrus.WithError(err).WithField("bot_id", agentID).Error("Cannot remove agent from LiveChat")
				}
			}(agent.Conn, agent.ID)
		}
	}

	m.apps = newApps
	wg.Wait()
	return nil
}

func (m *manager) Destroy(ctx context.Context) {
	wg := &sync.WaitGroup{}

	for _, app := range m.apps {
		wg.Add(1)

		go func(licenseID livechat.LicenseID) {
			logEntry := logrus.WithField("license_id", licenseID).WithContext(ctx)
			logEntry.Debug("Destroying app")
			err := m.UninstallApp(ctx, licenseID)
			if err != nil {
				logEntry.WithError(err).Error("Something went wrong during destroying the app")
			} else {
				logEntry.Debug("App destroyed")
			}
		}(app.licenseID)
	}

	wg.Wait()
}

func (m *manager) JoinChat(ctx context.Context, license livechat.LicenseID, chat livechat.ChatID) error {
	m.mu.Lock()
	var app *app
	for _, a := range m.apps {
		if a.licenseID == license {
			app = a
			break
		}
	}
	m.mu.Unlock()

	agents := []livechat.AgentID{}
	for _, agent := range app.bots {
		agents = append(agents, agent.ID)
	}

	return nil
}
