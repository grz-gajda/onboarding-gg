package rtm

import "github.com/livechat/onboarding/livechat"

type PushIncomingMessage struct {
	Action    string             `json:"action"`
	LicenseID livechat.LicenseID `json:"license_id,omitempty"`
	Payload   struct {
		ChatID   livechat.ChatID `json:"chat_id"`
		ThreadID string          `json:"thread_id,omitempty"`
		Event    struct {
			Type     string `json:"type"`
			Text     string `json:"text"`
			AuthorID string `json:"author_id"`
		} `json:"event"`
	} `json:"payload"`
}

func (m *PushIncomingMessage) GetAction() string                { return m.Action }
func (m *PushIncomingMessage) GetLicenseID() livechat.LicenseID { return m.LicenseID }

type PushIncomingChat struct {
	Action    string             `json:"action"`
	LicenseID livechat.LicenseID `json:"license_id,omitempty"`
	Payload   struct {
		Chat struct {
			ID livechat.ChatID `json:"id"`
		} `json:"chat"`
	} `json:"payload"`
}

func (m *PushIncomingChat) GetAction() string                { return m.Action }
func (m *PushIncomingChat) GetLicenseID() livechat.LicenseID { return m.LicenseID }
