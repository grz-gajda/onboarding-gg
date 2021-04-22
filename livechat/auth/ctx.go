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
	authorToken = ctxToken("with_author_token")
)

func WithPAT(ctx context.Context, username, password string) context.Context {
	token := fmt.Sprintf("Basic %s:%s", username, password)
	hashed := base64.StdEncoding.EncodeToString([]byte(token))

	return context.WithValue(ctx, authToken, hashed)
}

func WithOAuth(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authToken, fmt.Sprintf("Bearer %s", token))
}

func WithClientID(ctx context.Context, clientID livechat.ClientID) context.Context {
	return context.WithValue(ctx, clientToken, clientID)
}

func WithAuthorID(ctx context.Context, authorID livechat.AgentID) context.Context {
	return context.WithValue(ctx, authorToken, authorID)
}

func GetAuthToken(ctx context.Context) (string, error) {
	bearerToken, ok := ctx.Value(authToken).(string)
	if !ok {
		return "", fmt.Errorf("auth: missing authorization token")
	}
	if bearerToken == "" {
		return "", fmt.Errorf("auth: missing authorization token")
	}
	return bearerToken, nil
}

func GetClientID(ctx context.Context) (livechat.ClientID, error) {
	clientID, ok := ctx.Value(clientToken).(livechat.ClientID)
	if !ok {
		return "", fmt.Errorf("auth: missing client id")
	}
	if clientID == "" {
		return "", fmt.Errorf("auth: missing client id")
	}
	return clientID, nil
}

func GetAuthorID(ctx context.Context) (livechat.AgentID, error) {
	authorID, ok := ctx.Value(authorToken).(livechat.AgentID)
	if !ok {
		return "", fmt.Errorf("auth: missing author id")
	}
	if authorID == "" {
		return "", fmt.Errorf("auth: missing author id")
	}
	return authorID, nil
}
