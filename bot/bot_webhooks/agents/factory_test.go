package agents

import (
	"context"
	"errors"
	"testing"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/web/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Initialize(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	lcHTTP.On("CreateBot", ctx, mock.Anything).Once().Return(&livechat.CreateBotResponse{}, nil)

	lcHTTP.On("ListBots", ctx, mock.Anything).Once().Return([]*livechat.ListBotResponse{
		{ID: "abcd_1"},
		{ID: "abcd_2"},
	}, nil)

	lcHTTP.
		On("SetRoutingStatus", ctx, mock.MatchedBy(func(p *livechat.SetRoutingStatusRequest) bool { return p.AgentID == livechat.AgentID("abcd_1") })).
		Return(&livechat.SetRoutingStatusResponse{}, nil)

	lcHTTP.
		On("SetRoutingStatus", ctx, mock.MatchedBy(func(p *livechat.SetRoutingStatusRequest) bool { return p.AgentID == livechat.AgentID("abcd_2") })).
		Return(&livechat.SetRoutingStatusResponse{}, errors.New("invalid agent"))

	lcHTTP.
		On("DeleteBot", ctx, mock.MatchedBy(func(p *livechat.DeleteBotRequest) bool { return p.ID == livechat.AgentID("abcd_2") })).
		Once().
		Return(&livechat.DeleteBotResponse{}, nil)

	agents, err := Initialize(ctx, lcHTTP)
	assert.NoError(t, err)
	assert.Equal(t, 1, agents.Len())
}

func Test_Terminate(t *testing.T) {
	ctx := context.Background()
	lcHTTP := new(mocks.LivechatRequests)

	lcHTTP.On("SetRoutingStatus", ctx, mock.Anything).Return(&livechat.SetRoutingStatusResponse{}, nil)

	agents := NewCollection()
	agents.Register(NewAgent("abcd_1"))
	agents.Register(NewAgent("abcd_2"))
	agents.Register(NewAgent("abcd_3"))
	agents.Register(NewAgent("abcd_4"))

	err := Terminate(ctx, lcHTTP, agents)
	assert.NoError(t, err)
}

func Test_CreateBot(t *testing.T) {
	lcHTTP := new(mocks.LivechatRequests)
	ctx := context.Background()

	lcHTTP.
		On("CreateBot", ctx, mock.Anything).
		Return(&livechat.CreateBotResponse{ID: livechat.AgentID("abcd")}, nil)

	agent, err := createBot(ctx, lcHTTP)
	assert.NoError(t, err)
	assert.Equal(t, livechat.AgentID("abcd"), agent.ID)
}

func Test_FetchBots(t *testing.T) {
	lcHTTP := new(mocks.LivechatRequests)
	ctx := context.Background()

	lcHTTP.
		On("ListBots", ctx, mock.Anything).
		Return([]*livechat.ListBotResponse{
			{ID: "abcd_1"},
			{ID: "abcd_2"},
		}, nil)

	agents, err := fetchBots(ctx, lcHTTP)
	assert.NoError(t, err)
	assert.Equal(t, 2, agents.Len())
}

func Test_RemoveBot(t *testing.T) {
	lcHTTP := new(mocks.LivechatRequests)
	ctx := context.Background()

	lcHTTP.
		On("DeleteBot", ctx, mock.Anything).
		Return(&livechat.DeleteBotResponse{}, nil)

	err := removeBot(ctx, lcHTTP, livechat.AgentID("abcd_3"))
	assert.NoError(t, err)
}

func Test_EnableBot(t *testing.T) {
	lcHTTP := new(mocks.LivechatRequests)
	ctx := context.Background()

	lcHTTP.
		On("SetRoutingStatus", ctx, mock.MatchedBy(func(p *livechat.SetRoutingStatusRequest) bool {
			return p.Status == "accepting_chats"
		})).
		Return(&livechat.SetRoutingStatusResponse{}, nil)

	err := enableBot(ctx, lcHTTP, livechat.AgentID("abcd_3"))
	assert.NoError(t, err)
}

func Test_DisableBot(t *testing.T) {
	lcHTTP := new(mocks.LivechatRequests)
	ctx := context.Background()

	lcHTTP.
		On("SetRoutingStatus", ctx, mock.MatchedBy(func(p *livechat.SetRoutingStatusRequest) bool {
			return p.Status == "offline"
		})).
		Return(&livechat.SetRoutingStatusResponse{}, nil)

	err := disableBot(ctx, lcHTTP, livechat.AgentID("abcd_3"))
	assert.NoError(t, err)
}
