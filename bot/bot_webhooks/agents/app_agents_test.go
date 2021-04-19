package agents

import (
	"testing"

	"github.com/livechat/onboarding/livechat"
	"github.com/stretchr/testify/assert"
)

func Test_Agents(t *testing.T) {
	t.Run("register agent", func(t *testing.T) {
		agentsCollection := &collection{}
		agentsCollection.Register(&Agent{ID: "abcd"})

		assert.Len(t, agentsCollection.agents, 1)
	})

	t.Run("unregister agent", func(t *testing.T) {
		agentsCollection := &collection{}
		agentsCollection.Register(&Agent{ID: "abcd"})

		agent, err := agentsCollection.Unregister("abcd")
		assert.NoError(t, err)
		assert.Len(t, agentsCollection.agents, 0)
		assert.Equal(t, livechat.AgentID("abcd"), agent.ID)
	})

	t.Run("unregister agent (not found)", func(t *testing.T) {
		agentsCollection := &collection{}

		_, err := agentsCollection.Unregister("abcd")
		assert.Len(t, agentsCollection.agents, 0)
		assert.Error(t, err)
	})

	t.Run("find by chat exclude (success)", func(t *testing.T) {
		agentsCollection := &collection{}
		agentsCollection.Register(&Agent{ID: "abcd"})

		agent, err := agentsCollection.FindByChatExclude(livechat.ChatID("abcd"))
		assert.NoError(t, err)
		assert.Equal(t, livechat.AgentID("abcd"), agent.ID)
	})

	t.Run("find by chat exclude (fail)", func(t *testing.T) {
		agentsCollection := &collection{}

		_, err := agentsCollection.FindByChatExclude(livechat.ChatID("abcd"))
		assert.Error(t, err)
	})

	t.Run("find by chat (success)", func(t *testing.T) {
		agentsCollection := &collection{}
		agentsCollection.Register(&Agent{
			ID:    "abcd",
			chats: []livechat.ChatID{"abcd"},
		})

		agent, err := agentsCollection.FindByChat(livechat.ChatID("abcd"))
		assert.NoError(t, err)
		assert.Equal(t, livechat.AgentID("abcd"), agent.ID)
	})

	t.Run("find by chat (fail)", func(t *testing.T) {
		agentsCollection := &collection{}
		agentsCollection.Register(&Agent{ID: "abcd"})

		_, err := agentsCollection.FindByChat(livechat.ChatID("abcd"))
		assert.Error(t, err)
	})
}
