package auth

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/livechat/onboarding/livechat"
)

type ctxToken string

var (
	authToken   = ctxToken("auth_bearer")
	clientToken = ctxToken("client_id")
)

func WithToken(ctx context.Context, username, password string) context.Context {
	token := fmt.Sprintf("%s:%s", username, password)
	hashed := base64.StdEncoding.EncodeToString([]byte(token))

	return context.WithValue(ctx, authToken, hashed)
}

func WithClientID(ctx context.Context, clientID livechat.ClientID) context.Context {
	return context.WithValue(ctx, clientToken, clientID)
}

func GetAuthToken(ctx context.Context) (string, error) {
	bearerToken := ctx.Value(authToken).(string)
	if bearerToken == "" {
		return "", fmt.Errorf("auth: missing authorization token")
	}
	return bearerToken, nil
}

func GetClientID(ctx context.Context) (livechat.ClientID, error) {
	clientID := ctx.Value(clientToken).(livechat.ClientID)
	if clientID == "" {
		return "", fmt.Errorf("auth: missing client id")
	}
	return clientID, nil
}
