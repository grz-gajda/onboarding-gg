package bot_webhooks

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/livechat/onboarding/bot"
	"github.com/livechat/onboarding/bot/bot_webhooks/agents"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/web"
	log "github.com/sirupsen/logrus"
)

type manager struct {
	lcHTTP   web.LivechatRequests
	localURL string

	apps   *apps
	sender bot.Sender
}

func (m *manager) InstallApp(ctx context.Context, id livechat.LicenseID) error {
	app := newApp(m.lcHTTP, m.sender, id, m.localURL)
	m.apps.Register(app)

	bots, err := agents.Initialize(ctx, m.lcHTTP)
	if err != nil {
		return err
	}

	app.agents = bots

	if err := app.RegisterAction(ctx, WebhookEvents...); err != nil {
		log.WithField("license_id", id).WithError(err).Error("Cannot register 'incoming_chat' action")
		return err
	}
	if _, err := app.lcHTTP.EnableLicenseWebhook(ctx, &livechat.EnableLicenseWebhookRequest{}); err != nil {
		log.WithField("license_id", id).WithError(err).Error("Cannot enable webhooks")
		return err
	}

	return nil
}

func (m *manager) UninstallApp(ctx context.Context, id livechat.LicenseID) error {
	app := m.apps.Unregister(id)
	if app == nil {
		return fmt.Errorf("bot: app (license id: %v) is not registered", id)
	}

	if _, err := app.lcHTTP.DisableLicenseWebhook(ctx, &livechat.DisableLicenseWebhookRequest{}); err != nil {
		return err
	}

	return app.UnregisterActions(ctx)
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

func (m *manager) Redirect(ctx context.Context, rawMsg livechat.Push) error {
	app, err := m.apps.Find(rawMsg.GetLicenseID())
	if err != nil {
		return fmt.Errorf("bot: redirect_action: %w", err)
	}

	switch msg := rawMsg.(type) {
	case *livechat.PushIncomingMessage:
		log.Debug("Received *PushIncomingMessage")
		return app.IncomingEvent(ctx, msg)
	case *livechat.PushIncomingChat:
		log.Debug("Received *PushIncomingChat")
		return app.TransferChat(ctx, msg)
	}

	return errors.New("bot: received webhook with unknown message")
}
