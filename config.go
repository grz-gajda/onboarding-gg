package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/go-playground/validator"
	"github.com/livechat/onboarding/livechat"
)

type appMethod string

const (
	webhooksMethod = "webhooks"
	rtmMethod      = "rtm"
)

type config struct {
	Methods     appMethod   `json:"methods"`
	Auth        authConfig  `json:"auth" validate:"required"`
	Credentials credentials `json:"credentials" validate:"required"`
	URL         urlConfig   `json:"url" validate:"required"`
}

func (c *config) SelectMethod() appMethod {
	if c.Methods == "" {
		return webhooksMethod
	}
	return c.Methods
}

type authConfig struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type credentials struct {
	ClientID livechat.ClientID `json:"client_id" validate:"required"`
}

type urlConfig struct {
	HTTP string `json:"http" validate:"required"`
	WS   string `json:"ws" validate:"required"`
}

func LoadConfig(reader io.Reader) (*config, error) {
	var err error
	var cfg *config

	if err = json.NewDecoder(reader).Decode(&cfg); err != nil {
		return cfg, err
	}
	if err = validator.New().Struct(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func LoadConfigFile(filename string) (*config, error) {
	var err error
	var cfg *config

	file, err := os.Open(filename)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	cfg, err = LoadConfig(file)
	if err != nil {
		return cfg, err
	}

	return cfg, file.Close()
}
