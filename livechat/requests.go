package livechat

type Request interface {
	Endpoint() string
}
type ClientAuthorize interface {
	Request
	WithClientID(ClientID)
}

const (
	createBotEndpoint = "/configuration/action/create_bot"
	deleteBotEndpoint = "/configuration/action/delete_bot"
	listBotsEndpoint  = "/configuration/action/list_bots"

	transferChatEndpoint = "/agent/action/transfer_chat"
	sendEventEndpoint    = "/agent/action/send_event"

	registerWebhookEndpoint       = "/configuration/action/register_webhook"
	unregisterWebhookEndpoint     = "/configuration/action/unregister_webhook"
	enableLicenseWebhookEndpoint  = "/configuration/action/enable_license_webhooks"
	disableLicenseWebhookEndpoint = "/configuration/action/disable_license_webhooks"

	setRoutingStatusEndpoint = "/agent/action/set_routing_status"
)

type CreateBotRequest struct {
	Name     string   `json:"name"`
	ClientID ClientID `json:"owner_client_id,omitempty"`
}

func (r *CreateBotRequest) Endpoint() string          { return createBotEndpoint }
func (r *CreateBotRequest) WithClientID(cid ClientID) { r.ClientID = cid }

type CreateBotResponse struct {
	ID AgentID `json:"id"`
}

type DeleteBotRequest struct {
	ID AgentID `json:"id"`
}

func (r *DeleteBotRequest) Endpoint() string { return deleteBotEndpoint }

type DeleteBotResponse struct{}

type ListBotsRequest struct {
	All bool `json:"all,omitempty"`
}

func (r *ListBotsRequest) Endpoint() string { return listBotsEndpoint }

type ListBotResponse struct {
	ID   AgentID `json:"id"`
	Name string  `json:"name"`
}

type RegisterWebhookRequest struct {
	Action    string   `json:"action"`
	SecretKey string   `json:"secret_key"`
	URL       string   `json:"url"`
	Type      string   `json:"type"`
	ClientID  ClientID `json:"owner_client_id,omitempty"`
}

func (r *RegisterWebhookRequest) Endpoint() string          { return registerWebhookEndpoint }
func (r *RegisterWebhookRequest) WithClientID(cid ClientID) { r.ClientID = cid }

type RegisterWebhookResponse struct {
	ID string `json:"id"`
}

type UnregisterWebhookRequest struct {
	ID       string   `json:"id"`
	ClientID ClientID `json:"owner_client_id,omitempty"`
}

func (r *UnregisterWebhookRequest) Endpoint() string          { return unregisterWebhookEndpoint }
func (r *UnregisterWebhookRequest) WithClientID(cid ClientID) { r.ClientID = cid }

type UnregisterWebhookResponse struct{}

type TransferChatRequest struct {
	ID     ChatID `json:"id"`
	Target struct {
		Type string    `json:"type"`
		IDs  []AgentID `json:"ids"`
	} `json:"target"`
	Force bool `json:"force,omitempty"`
}

func (r *TransferChatRequest) Endpoint() string { return transferChatEndpoint }

type TransferChatResponse struct{}

type SendEventResponse struct {
	EventID string `json:"event_id"`
}

type EnableLicenseWebhookRequest struct {
	ClientID ClientID `json:"owner_client_id,omitempty"`
}

func (r *EnableLicenseWebhookRequest) Endpoint() string          { return enableLicenseWebhookEndpoint }
func (r *EnableLicenseWebhookRequest) WithClientID(cid ClientID) { r.ClientID = cid }

type EnableLicenseWebhookResponse struct{}

type DisableLicenseWebhookRequest struct {
	ClientID ClientID `json:"owner_client_id,omitempty"`
}

func (r *DisableLicenseWebhookRequest) Endpoint() string          { return disableLicenseWebhookEndpoint }
func (r *DisableLicenseWebhookRequest) WithClientID(cid ClientID) { r.ClientID = cid }

type DisableLicenseWebhookResponse struct{}

type SetRoutingStatusRequest struct {
	Status  string  `json:"status"`
	AgentID AgentID `json:"agent_id"`
}

func (r *SetRoutingStatusRequest) Endpoint() string { return setRoutingStatusEndpoint }

type SetRoutingStatusResponse struct{}
