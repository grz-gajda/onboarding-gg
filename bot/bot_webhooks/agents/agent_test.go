package agents

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Agent_RegisterChat(t *testing.T) {
	agent := NewAgent("abcd_1")

	agent.RegisterChat("chat_1", "chat_2", "chat_3")
	assert.Len(t, agent.chats, 3)

	agent.RegisterChat("chat_1", "chat_4", "chat_5")
	assert.Len(t, agent.chats, 5)
}
