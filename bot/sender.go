package bot

import (
	"context"
	"strings"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/web"
	log "github.com/sirupsen/logrus"
)

type sender struct {
	client      web.LivechatRequests
	appAuthorID string
}

func NewSender(client web.LivechatRequests, authorID string) Sender {
	return &sender{client: client, appAuthorID: authorID}
}

func (s *sender) Talk(ctx context.Context, chatID livechat.ChatID, msg *livechat.PushIncomingMessage) error {
	if msg.Payload.Event.Type != "message" {
		return nil
	}
	if msg.Payload.Event.AuthorID == s.appAuthorID {
		return nil
	}

	switch msg.Payload.Event.Text {
	case "Hello":
		_, err := s.client.SendEvent(ctx, livechat.BuildMessage(chatID, "World!"))
		return err
	case "Wróć do człowieka":
		return s.redirectToAgent(ctx, chatID)
	default:
		_, err := s.client.SendEvent(ctx, livechat.BuildButtonMessage(chatID, "Czy chcesz wrócić do człowieka?", "", "Wróć do człowieka"))
		return err
	}
}

func (s *sender) redirectToAgent(ctx context.Context, chatID livechat.ChatID) error {
	realAgents, err := s.client.ListAgentsForTransfer(ctx, &livechat.ListAgentsForTransferRequest{ChatID: chatID})
	if err != nil {
		log.WithError(err).Error("Cannot fetch list of real agents")
		return err
	}

	if len(realAgents) == 0 {
		_, err = s.client.SendEvent(ctx, livechat.BuildMessage(chatID, "Obecnie nie ma żadnego człowieka do rozmowy :("))
		return err
	}

	for _, realAgent := range realAgents {
		if _, err = s.client.TransferChat(ctx, buildTransferChatMessage(chatID, realAgent.AgentID)); err != nil {
			if isAgentOffline(err) || isAgentAssigned(err) {
				continue
			}

			log.WithError(err).WithField("chat_id", chatID).Error("Cannot transfer chat to agent on demand")
			continue
		}
		return nil
	}

	_, err = s.client.SendEvent(ctx, livechat.BuildMessage(chatID, "Obecnie nie ma żadnego człowieka do rozmowy :("))
	return err
}

func buildTransferChatMessage(chatID livechat.ChatID, agentID ...livechat.AgentID) *livechat.TransferChatRequest {
	return &livechat.TransferChatRequest{
		ID: chatID,
		Target: struct {
			Type string             "json:\"type\""
			IDs  []livechat.AgentID "json:\"ids\""
		}{
			Type: "agent",
			IDs:  agentID,
		},
		Force: false,
	}
}

func isAgentOffline(err error) bool {
	return strings.Contains(err.Error(), "Agent is offline.")
}
func isAgentAssigned(err error) bool {
	return strings.Contains(err.Error(), "One or more of requested agents are already present in the chat")
}
