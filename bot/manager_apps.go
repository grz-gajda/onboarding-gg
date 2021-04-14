package bot

import (
	"fmt"
	"sync"

	"github.com/livechat/onboarding/livechat"
)

type apps struct {
	apps []*app
	mu   sync.Mutex
}

func (a *apps) Register(app *app) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, registeredApp := range a.apps {
		if registeredApp.licenseID == app.licenseID {
			return fmt.Errorf("bot: app (license id: %v) is already installed", app.licenseID)
		}
	}

	a.apps = append(a.apps, app)
	return nil
}

func (a *apps) Unregister(licenseID livechat.LicenseID) *app {
	a.mu.Lock()
	defer a.mu.Unlock()

	newApps := []*app{}
	var uninstalledApp *app

	for _, registeredApp := range a.apps {
		if registeredApp.licenseID != licenseID {
			newApps = append(newApps, registeredApp)
		} else {
			uninstalledApp = registeredApp
		}
	}

	a.apps = newApps
	return uninstalledApp
}

func (a *apps) Find(licenseID livechat.LicenseID) (*app, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, app := range a.apps {
		if app.licenseID == licenseID {
			return app, nil
		}
	}
	return nil, fmt.Errorf("bot: app (license id: %v) is not installed", licenseID)
}
