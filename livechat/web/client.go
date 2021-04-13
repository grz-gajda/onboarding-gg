package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/livechat/onboarding/livechat/auth"
)

type livechatClient struct {
	httpClient Client
	url        string
}

type errorResponse struct {
	ErrorMessage struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func (c *livechatClient) sendRequest(ctx context.Context, endpoint string, payload interface{}, body interface{}) (*http.Response, error) {
	var err error
	var req *http.Request

	jsonBody := []byte("{}")
	if payload != nil {
		jsonBody, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("http_client: cannot encode request body: %w", err)
		}
	}

	url := fmt.Sprintf("%s%s", c.url, endpoint)
	req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("http_client: %w", err)
	}

	bearerToken, err := auth.GetAuthToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("http_client: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", bearerToken))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http_client: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, readErrorMessage(res.Body)
	}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return res, fmt.Errorf("http_client: %w", err)
	}
	return res, nil
}

func readErrorMessage(w io.Reader) error {
	body := errorResponse{}
	if err := json.NewDecoder(w).Decode(&body); err != nil {
		return fmt.Errorf("http_client: cannot decode response: %w", err)
	}

	return fmt.Errorf("%s (type %s)", body.ErrorMessage.Message, body.ErrorMessage.Type)
}
