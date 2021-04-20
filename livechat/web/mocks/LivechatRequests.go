// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	context "context"

	livechat "github.com/livechat/onboarding/livechat"
	mock "github.com/stretchr/testify/mock"
)

// LivechatRequests is an autogenerated mock type for the LivechatRequests type
type LivechatRequests struct {
	mock.Mock
}

// CreateBot provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) CreateBot(_a0 context.Context, _a1 *livechat.CreateBotRequest) (*livechat.CreateBotResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *livechat.CreateBotResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.CreateBotRequest) *livechat.CreateBotResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*livechat.CreateBotResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.CreateBotRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteBot provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) DeleteBot(_a0 context.Context, _a1 *livechat.DeleteBotRequest) (*livechat.DeleteBotResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *livechat.DeleteBotResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.DeleteBotRequest) *livechat.DeleteBotResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*livechat.DeleteBotResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.DeleteBotRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DisableLicenseWebhook provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) DisableLicenseWebhook(_a0 context.Context, _a1 *livechat.DisableLicenseWebhookRequest) (*livechat.DisableLicenseWebhookResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *livechat.DisableLicenseWebhookResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.DisableLicenseWebhookRequest) *livechat.DisableLicenseWebhookResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*livechat.DisableLicenseWebhookResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.DisableLicenseWebhookRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EnableLicenseWebhook provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) EnableLicenseWebhook(_a0 context.Context, _a1 *livechat.EnableLicenseWebhookRequest) (*livechat.EnableLicenseWebhookResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *livechat.EnableLicenseWebhookResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.EnableLicenseWebhookRequest) *livechat.EnableLicenseWebhookResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*livechat.EnableLicenseWebhookResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.EnableLicenseWebhookRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAgents provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) ListAgents(_a0 context.Context, _a1 *livechat.ListAgentsRequest) ([]*livechat.ListAgentsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*livechat.ListAgentsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.ListAgentsRequest) []*livechat.ListAgentsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*livechat.ListAgentsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.ListAgentsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListBots provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) ListBots(_a0 context.Context, _a1 *livechat.ListBotsRequest) ([]*livechat.ListBotResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*livechat.ListBotResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.ListBotsRequest) []*livechat.ListBotResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*livechat.ListBotResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.ListBotsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterWebhook provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) RegisterWebhook(_a0 context.Context, _a1 *livechat.RegisterWebhookRequest) (*livechat.RegisterWebhookResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *livechat.RegisterWebhookResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.RegisterWebhookRequest) *livechat.RegisterWebhookResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*livechat.RegisterWebhookResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.RegisterWebhookRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendEvent provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) SendEvent(_a0 context.Context, _a1 *livechat.Event) (*livechat.SendEventResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *livechat.SendEventResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.Event) *livechat.SendEventResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*livechat.SendEventResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.Event) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetRoutingStatus provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) SetRoutingStatus(_a0 context.Context, _a1 *livechat.SetRoutingStatusRequest) (*livechat.SetRoutingStatusResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *livechat.SetRoutingStatusResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.SetRoutingStatusRequest) *livechat.SetRoutingStatusResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*livechat.SetRoutingStatusResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.SetRoutingStatusRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransferChat provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) TransferChat(_a0 context.Context, _a1 *livechat.TransferChatRequest) (*livechat.TransferChatResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *livechat.TransferChatResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.TransferChatRequest) *livechat.TransferChatResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*livechat.TransferChatResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.TransferChatRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnregisterWebhook provides a mock function with given fields: _a0, _a1
func (_m *LivechatRequests) UnregisterWebhook(_a0 context.Context, _a1 *livechat.UnregisterWebhookRequest) (*livechat.UnregisterWebhookResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *livechat.UnregisterWebhookResponse
	if rf, ok := ret.Get(0).(func(context.Context, *livechat.UnregisterWebhookRequest) *livechat.UnregisterWebhookResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*livechat.UnregisterWebhookResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *livechat.UnregisterWebhookRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
