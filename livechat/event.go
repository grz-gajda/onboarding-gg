package livechat

type EventType string
type ButtonType string

const (
	EventTypeMessage     EventType = "message"
	EventTypeRichMessage EventType = "rich_message"
)

const (
	ButtonTypeWebView ButtonType = "webview"
	ButtonTypeMessage ButtonType = "message"
	ButtonTypeUrl     ButtonType = "url"
	ButtonTypePhone   ButtonType = "phone"
)

type Event struct {
	ChatID ChatID       `json:"chat_id"`
	Event  EventMessage `json:"event"`
}

func (r *Event) Endpoint() string       { return sendEventEndpoint }
func (r *Event) RequiresClientID() bool { return false }

type EventMessage struct {
	Type     EventType      `json:"type"`
	Text     string         `json:"text,omitempty"`
	Elements []EventElement `json:"elements,omitempty"`
}

type EventElement struct {
	Title    string        `json:"title"`
	SubTitle string        `json:"subtitle"`
	Image    string        `json:"image"`
	Button   []EventButton `json:"buttons,omitempty"`
}

type EventButton struct {
	Text string     `json:"text"`
	Type ButtonType `json:"type"`
}
