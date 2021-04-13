package web

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/livechat/onboarding/livechat/auth"
	"github.com/livechat/onboarding/livechat/web/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Client_WithoutBody(t *testing.T) {
	type exampleResponse struct {
		Body string `json:"body"`
	}

	ctx := auth.WithToken(context.Background(), "username", "password")

	httpClient := new(mocks.Client)
	httpClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"body": "example response"}`)),
	}, nil)

	webService := &livechatClient{
		httpClient: httpClient,
		url:        "http://lorem.pl",
	}

	var example exampleResponse
	_, err := webService.sendRequest(ctx, "/random/endpoint", nil, &example)

	assert.NoError(t, err)
	assert.Equal(t, "example response", example.Body)
}

func Test_Client_WithBody(t *testing.T) {
	type exampleRequest struct {
		Data string `json:"data"`
	}
	type exampleResponse struct {
		Body string `json:"body"`
	}

	ctx := auth.WithToken(context.Background(), "username", "password")

	httpClient := new(mocks.Client)
	httpClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		var b exampleRequest
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			return false
		}

		return b.Data == "random data"
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"body": "example response"}`)),
	}, nil)

	webService := &livechatClient{
		httpClient: httpClient,
		url:        "http://lorem.pl",
	}

	var example exampleResponse
	_, err := webService.sendRequest(ctx, "/random/endpoint", exampleRequest{Data: "random data"}, &example)

	assert.NoError(t, err)
	assert.Equal(t, "example response", example.Body)
}
