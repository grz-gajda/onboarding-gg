package livechat

import "net/http"

type LicenseID int
type ChatID string
type ClientID string
type AgentID string

//go:generate mockery --name Client
type Client interface {
	Do(*http.Request) (*http.Response, error)
}
