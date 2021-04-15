package bot_webhooks

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
	log "github.com/sirupsen/logrus"
)

type manager struct {
	lcHTTP   web.LivechatRequests
	localURL string

	apps       *apps
	botFactory *botFactory
}

func (m *manager) InstallApp(ctx context.Context, id livechat.LicenseID) error {
	app := newApp(m.lcHTTP, id)
	m.apps.Register(app)

	clientID, err := auth.GetClientID(ctx)
	if err != nil {
		return err
	}

	bots, err := m.botFactory.Initialize(ctx)
	if err != nil {
		return err
	}
	for _, bot := range bots {
		app.agents.Register(bot)
	}

	if err := app.RegisterAction(ctx, "incoming_chat", &registerActionOptions{CustomURL: m.localURL}); err != nil {
		log.WithField("license_id", id).WithError(err).Error("Cannot register 'incoming_chat' action")
		return err
	}
	if err := app.RegisterAction(ctx, "incoming_event", &registerActionOptions{CustomURL: m.localURL}); err != nil {
		log.WithField("license_id", id).WithError(err).Error("Cannot register 'incoming_event' action")
		return err
	}
	if _, err := app.lcHTTP.EnableLicenseWebhook(ctx, &web.EnableLicenseWebhookRequest{ClientID: clientID}); err != nil {
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

	clientID, err := auth.GetClientID(ctx)
	if err != nil {
		return err
	}

	if _, err := app.lcHTTP.DisableLicenseWebhook(ctx, &web.DisableLicenseWebhookRequest{ClientID: clientID}); err != nil {
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

func (m *manager) Redirect(ctx context.Context, rawMsg rtm.Push, payload ...RedirectData) error {
	app, err := m.apps.Find(rawMsg.GetLicenseID())
	if err != nil {
		return fmt.Errorf("bot: redirect_action: %w", err)
	}

	switch msg := rawMsg.(type) {
	case *rtm.PushIncomingMessage:
		log.WithField("event", msg).Debug("Received *PushIncomingMessage")
		return app.IncomingEvent(ctx, msg, payload...)
	case *rtm.PushIncomingChat:
		log.WithField("event", msg).Debug("Received *PushIncomingChat")
		return app.TransferChat(ctx, msg, payload...)
	}

	return errors.New("bot: received webhook with unknown message")
}
