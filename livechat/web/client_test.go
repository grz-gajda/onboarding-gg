package web

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/livechat/onboarding/livechat"
	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/web/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type exampleRequest struct {
	Data string `json:"data"`
}

func (r *exampleRequest) Endpoint() string { return "/random/endpoint" }

type exampleRequestClientID struct {
	Data     string            `json:"data"`
	ClientID livechat.ClientID `json:"owner_client_id,omitempty"`
}

func (r *exampleRequestClientID) Endpoint() string                   { return "/random/example_request" }
func (r *exampleRequestClientID) WithClientID(cid livechat.ClientID) { r.ClientID = cid }

type exampleResponse struct {
	Body string `json:"body"`
}

func Test_Client_WithBody(t *testing.T) {
	ctx := auth.WithToken(context.Background(), "username", "password")

	httpClient := new(mocks.Client)
	httpClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		var b exampleRequest
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			return false
		}

		return b.Data == "random data" && r.URL.Path == "/random/endpoint"
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"body": "example response"}`)),
	}, nil)

	webService := &livechatClient{
		httpClient: httpClient,
		url:        "http://lorem.pl",
	}

	var example exampleResponse
	_, err := webService.sendRequest(ctx, &exampleRequest{Data: "random data"}, &example)

	assert.NoError(t, err)
	assert.Equal(t, "example response", example.Body)
}

func Test_Client_WithClientID(t *testing.T) {
	ctx := auth.WithToken(context.Background(), "username", "password")
	ctx = auth.WithClientID(ctx, livechat.ClientID("custom_client_id"))

	httpClient := new(mocks.Client)
	httpClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		var b exampleRequestClientID
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			return false
		}

		return b.ClientID == livechat.ClientID("custom_client_id") && b.Data == "random data" && r.URL.Path == "/random/example_request"
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"body": "example response"}`)),
	}, nil)

	webService := &livechatClient{
		httpClient: httpClient,
		url:        "http://lorem.pl",
	}

	var example exampleResponse
	_, err := webService.sendRequest(ctx, &exampleRequestClientID{Data: "random data"}, &example)

	assert.NoError(t, err)
	assert.Equal(t, "example response", example.Body)
}
