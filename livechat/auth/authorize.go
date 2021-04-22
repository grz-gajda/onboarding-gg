package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/livechat/onboarding/livechat"
	"github.com/sirupsen/logrus"
)

type AuthorizeCredentials struct {
	Code        string
	ClientID    livechat.ClientID
	Secret      string
	RedirectURI string
}

type AuthorizationResponse struct {
	AccessToken string `json:"access_token"`
}

type authErrorMessage struct {
	Error string `json:"error"`
	Desc  string `json:"error_description"`
}

func Authorize(ctx context.Context, client livechat.Client, data *AuthorizeCredentials) (*AuthorizationResponse, error) {
	b := url.Values{}
	b.Set("grant_type", "authorization_code")
	b.Set("code", data.Code)
	b.Set("client_id", string(data.ClientID))
	b.Set("client_secret", data.Secret)
	b.Set("redirect_uri", data.RedirectURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://accounts.livechat.com/v2/token", strings.NewReader(b.Encode()))
	if err != nil {
		return &AuthorizationResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(b.Encode())))

	res, err := client.Do(req)
	if err != nil {
		return &AuthorizationResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var body authErrorMessage
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return &AuthorizationResponse{}, err
		}

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"status_code": res.StatusCode,
			"error_type":  body.Error,
		}).Error(body.Desc)

		return &AuthorizationResponse{}, fmt.Errorf("auth: expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	var token AuthorizationResponse
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return &AuthorizationResponse{}, err
	}

	return &token, nil
}
