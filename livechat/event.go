package livechat

type EventType string
type ButtonType string
type TemplateID string

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

const (
	TemplateCards      TemplateID = "cards"
	TemplateSticker    TemplateID = "sticker"
	TemplateQuickReply TemplateID = "quick_replies"
)

type Event struct {
	ChatID ChatID       `json:"chat_id"`
	Event  EventMessage `json:"event"`
}

func (r *Event) Endpoint() string       { return sendEventEndpoint }
func (r *Event) RequiresClientID() bool { return false }

type EventMessage struct {
	Type       EventType      `json:"type"`
	Text       string         `json:"text,omitempty"`
	TemplateID TemplateID     `json:"template_id,omitempty"`
	Elements   []EventElement `json:"elements,omitempty"`
}

type EventElement struct {
	Title    string        `json:"title"`
	SubTitle string        `json:"subtitle"`
	Image    EventImage    `json:"image,omitempty"`
	Button   []EventButton `json:"buttons,omitempty"`
}

type EventImage struct {
	URL string `json:"url"`
}

type EventButton struct {
	Text       string     `json:"text"`
	Type       ButtonType `json:"type"`
	Value      string     `json:"value"`
	PostbackID string     `json:"postback_id"`
	UserID     []string   `json:"user_ids"`
}

func BuildMessage(chatID ChatID, text string) *Event {
	return &Event{
		ChatID: chatID,
		Event: EventMessage{
			Type: EventTypeMessage,
			Text: text,
		},
	}
}

func BuildButtonMessage(chatID ChatID, title, subtitle, text string) *Event {
	return &Event{
		ChatID: chatID,
		Event: EventMessage{
			Type:       EventTypeRichMessage,
			TemplateID: TemplateCards,
			Elements: []EventElement{{
				Title:    title,
				SubTitle: subtitle,
				Image: EventImage{
					URL: "https://en.meming.world/images/en/thumb/2/2c/Surprised_Pikachu_HD.jpg/300px-Surprised_Pikachu_HD.jpg",
				},
				Button: []EventButton{{
					Text:       text,
					Type:       ButtonTypeMessage,
					Value:      text,
					PostbackID: "to tez nie wiem co to za bardzo",
					UserID:     []string{},
				}},
			}},
		},
	}
}
