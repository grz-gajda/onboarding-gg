package bot_webhooks

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/livechat/onboarding/bot"
	"github.com/livechat/onboarding/bot/bot_webhooks/agents"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/web"
	log "github.com/sirupsen/logrus"
)

type manager struct {
	lcHTTP   web.LivechatRequests
	localURL string

	apps   *apps
	sender bot.Sender

	muAuth         *sync.Mutex
	authToken      string
	readyToInstall chan bool
}

func (m *manager) Authorize(ctx context.Context, client livechat.Client, data *auth.AuthorizeCredentials) error {
	m.muAuth.Lock()
	defer m.muAuth.Unlock()

	if m.authToken != "" {
		m.authToken = ""
		m.readyToInstall = make(chan bool, 1)
	}

	response, err := auth.Authorize(ctx, client, data)
	if err != nil {
		return err
	}

	m.authToken = response.AccessToken
	m.readyToInstall <- true
	close(m.readyToInstall)

	return nil
}

func (m *manager) InstallApp(ctx context.Context, id livechat.LicenseID) error {
	app := newApp(m.lcHTTP, m.sender, id, m.localURL)
	m.apps.Register(app)

	if m.authToken == "" {
		select {
		case <-m.readyToInstall:
			log.Debug("App is ready to be installed!")
			break
		case <-time.After(30 * time.Second):
			return errors.New("timeout")
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	ctx = auth.WithOAuth(ctx, m.authToken)
	bots, err := agents.Initialize(ctx, m.lcHTTP)
	if err != nil {
		return err
	}

	app.agents = bots

	if err := app.RegisterAction(ctx, WebhookEvents...); err != nil {
		log.WithField("license_id", id).WithError(err).Error("Cannot register 'incoming_chat' action")
		m.apps.Unregister(id)
		return err
	}
	if _, err := app.lcHTTP.EnableLicenseWebhook(ctx, &livechat.EnableLicenseWebhookRequest{}); err != nil {
		log.WithField("license_id", id).WithError(err).Error("Cannot enable webhooks")
		m.apps.Unregister(id)
		return err
	}

	return nil
}

func (m *manager) UninstallApp(ctx context.Context, id livechat.LicenseID) error {
	app := m.apps.Unregister(id)
	if app == nil {
		return fmt.Errorf("bot: app (license id: %v) is not registered", id)
	}

	defer func() {
		m.authToken = ""
		m.readyToInstall = make(chan bool, 1)
	}()

	ctx = auth.WithOAuth(ctx, m.authToken)
	if _, err := app.lcHTTP.DisableLicenseWebhook(ctx, &livechat.DisableLicenseWebhookRequest{}); err != nil {
		return err
	}

	return app.UnregisterActions(ctx)
}

func (m *manager) Destroy(ctx context.Context) {
	wg := &sync.WaitGroup{}
	ctx = auth.WithOAuth(ctx, m.authToken)

	for _, app := range m.apps.apps {
		wg.Add(1)

		go func(id livechat.LicenseID) {
			defer wg.Done()
			m.UninstallApp(ctx, id)
		}(app.licenseID)
	}

	wg.Wait()
}

func (m *manager) Redirect(ctx context.Context, rawMsg livechat.Push) error {
	app, err := m.apps.Find(rawMsg.GetLicenseID())
	if err != nil {
		return fmt.Errorf("bot: redirect_action: %w", err)
	}

	ctx = auth.WithOAuth(ctx, m.authToken)
	logEntry := log.WithFields(log.Fields{
		"license_id": rawMsg.GetLicenseID(),
		"action":     rawMsg.GetAction(),
	})

	switch msg := rawMsg.(type) {
	case *livechat.PushIncomingMessage:
		logEntry.Debug("Received *PushIncomingMessage")
		return app.IncomingEvent(ctx, msg)
	case *livechat.PushIncomingChat:
		logEntry.Debug("Received *PushIncomingChat")
		return app.TransferChat(ctx, msg)
	case *livechat.PushUserAddedToChat:
		logEntry.WithField("raw_message", rawMsg).Debug("Received *PushUserAddedToChat")
		return app.UserAddedToChat(ctx, msg)
	default:
		logEntry.Warn("Received webhook with unknown message")
		return errors.New("bot: received webhook with unknown message")
	}
}
