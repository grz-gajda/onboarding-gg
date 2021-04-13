package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/livechat/onboarding/bot"
	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/rtm"
	"github.com/livechat/onboarding/livechat/web"
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

	// TOOLS
	httpClient := &http.Client{Timeout: 5 * time.Second}

	// GLOBAL CONTEXT
	ctx := auth.WithToken(context.Background(), cfg.Auth.Username, cfg.Auth.Password)
	ctx = auth.WithClientID(ctx, cfg.Credentials.ClientID)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// LIVECHAT SERVICES
	lcHTTP := web.New(httpClient, cfg.URL.HTTP)
	lcRTM, err := rtm.New(ctx, &websocket.Dialer{HandshakeTimeout: 5 * time.Second}, cfg.URL.WS)
	if err != nil {
		log.WithError(err).Panic("Cannot initialize connection to LiveChat")
	}

	if err := lcRTM.Login(ctx); err != nil {
		log.WithError(err).Panic("Cannot authorize connection to LiveChat")
	}

	botManager := bot.New(lcHTTP, lcRTM)

	// BACKGROUND
	go lcRTM.Ping(ctx)
	Shutdown(ctx, cancel, func() {
		httpClient.CloseIdleConnections()
		botManager.Destroy(ctx)
	})

	// HTTP
	router := chi.NewRouter()
	router.Use(middleware.RequestLogger(&logrusFormatter{logger: log.StandardLogger()}))
	router.Use(middleware.Recoverer)

	router.Post("/webhooks/install", func(w http.ResponseWriter, r *http.Request) {
		var payload webhooksPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			sendError(w, err)
			return
		}

		if payload.Event == "application_installed" {
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

	router.Post("/webhooks/incoming_chat", func(w http.ResponseWriter, r *http.Request) {
		var payload incomingChatPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			sendError(w, err)
			return
		}

		if payload.Action != "incoming_chat" {
			sendError(w, errors.New("received unknown action"))
		}

		if err := botManager.JoinChat(ctx, payload.LicenseID, payload.Payload.Chat.ID); err != nil {
			sendError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Print("Starting application")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Panicf("Something happened during HTTP request: %s", err)
	}
}

func sendError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": err.Error(),
	})
}

type webhooksPayload struct {
	LicenseID livechat.LicenseID `json:"licenseID"`
	AppName   string             `json:"applicationName"`
	ClientID  string             `json:"clientID"`
	Event     string             `json:"event"`
}

type incomingChatPayload struct {
	Action    string             `json:"action"`
	LicenseID livechat.LicenseID `json:"license_id"`
	Payload   struct {
		Chat struct {
			ID livechat.ChatID `json:"id"`
		} `json:"chat"`
	} `json:"payload"`
}
