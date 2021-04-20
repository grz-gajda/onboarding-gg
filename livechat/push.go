package livechat

type Push interface {
	GetAction() string
	GetLicenseID() LicenseID
}

type InstallApplicationWebhook struct {
	LicenseID LicenseID `json:"licenseID"`
	AppName   string    `json:"applicationName"`
	ClientID  string    `json:"clientID"`
	Event     string    `json:"event"`
}

type PushIncomingMessage struct {
	Action    string    `json:"action"`
	LicenseID LicenseID `json:"license_id,omitempty"`
	Payload   struct {
		ChatID   ChatID `json:"chat_id"`
		ThreadID string `json:"thread_id,omitempty"`
		Event    struct {
			Type     string `json:"type"`
			Text     string `json:"text"`
			AuthorID string `json:"author_id"`
		} `json:"event"`
	} `json:"payload"`
}

func (m *PushIncomingMessage) GetAction() string       { return m.Action }
func (m *PushIncomingMessage) GetLicenseID() LicenseID { return m.LicenseID }

type PushIncomingChat struct {
	Action    string    `json:"action"`
	LicenseID LicenseID `json:"license_id,omitempty"`
	Payload   struct {
		Chat struct {
			ID ChatID `json:"id"`
		} `json:"chat"`
	} `json:"payload"`
}

func (m *PushIncomingChat) GetAction() string       { return m.Action }
func (m *PushIncomingChat) GetLicenseID() LicenseID { return m.LicenseID }

type PushUserAddedToChat struct {
	Action    string    `json:"action"`
	LicenseID LicenseID `json:"license_id,omitempty"`
	Payload   struct {
		ChatID ChatID `json:"chat_id"`
		User   struct {
			Present bool   `json:"present"`
			Type    string `json:"type"`
		} `json:"user"`
	} `json:"payload"`
}

func (m *PushUserAddedToChat) GetAction() string       { return m.Action }
func (m *PushUserAddedToChat) GetLicenseID() LicenseID { return m.LicenseID }
