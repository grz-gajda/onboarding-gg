package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	// CONFIGURATION
	cfg, err := LoadConfigFile("config.json")
	if err != nil {
		log.WithError(err).Panic("Cannot load configuration for app")
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}

	// GLOBAL CONTEXT
	ctx := auth.WithClientID(context.Background(), cfg.Credentials.ClientID)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	router := chi.NewRouter()
	router.Use(middleware.RequestLogger(&logrusFormatter{logger: log.StandardLogger()}))
	router.Use(middleware.Recoverer)

	botManager := StartWebhooks(cfg, &appMethodConfig{
		httpClient: httpClient,
		router:     router,
	})

	Shutdown(ctx, cancel, func() {
		httpClient.CloseIdleConnections()
		botManager.Destroy(ctx)
	})

	router.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		referrerUri := fmt.Sprintf("%s/auth", cfg.URL.Local)

		code := r.URL.Query().Get("code")
		if code == "" {
			w.Header().Add("Location", fmt.Sprintf("https://accounts.livechat.com/?response_type=code&client_id=%s&redirect_uri=%s", cfg.Credentials.ClientID, referrerUri))
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		err := botManager.Authorize(r.Context(), httpClient, &auth.AuthorizeCredentials{
			Code:        code,
			ClientID:    cfg.Credentials.ClientID,
			Secret:      cfg.Credentials.Secret,
			RedirectURI: referrerUri,
		})

		if err != nil {
			sendError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	router.Post("/webhooks/install", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload livechat.InstallApplicationWebhook
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			sendError(w, err)
			return
		}

		if payload.Event == "application_installed" {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := botManager.InstallApp(ctx, payload.LicenseID); err != nil {
				log.WithError(err).WithField("id", payload.LicenseID).Error("Cannot install application")
				sendError(w, err)
				return
			}
		}
		if payload.Event == "application_uninstalled" {
			if err := botManager.UninstallApp(ctx, payload.LicenseID); err != nil {
				log.WithError(err).WithField("id", payload.LicenseID).Error("Cannot uninstall application")
				sendError(w, err)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Print("Starting application")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Panicf("Something happened during HTTP request: %s", err)
	}
}

func sendError(w http.ResponseWriter, err error) {
	if err != nil {
		log.WithError(err).Error("Outcoming invalid HTTP response")
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": err.Error(),
	})
}

type appMethodConfig struct {
	httpClient *http.Client
	router     *chi.Mux
}
