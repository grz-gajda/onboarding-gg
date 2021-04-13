package main

import (
	"bytes"
	"testing"
)

func Test_LoadConfig_Valid(t *testing.T) {
	_, err := LoadConfigFile("./config.dist.json")
	if err != nil {
		t.Fatalf("LoadConfigFile returns non-empty err: %s", err)
	}
}

func Test_LoadConfig_Invalid(t *testing.T) {
	_, err := LoadConfigFile("./config.never.json")
	if err == nil {
		t.Fatalf("LoadConfigFile returns empty err")
	}
}

func Test_LoadConfig_InvalidContent(t *testing.T) {
	content := bytes.NewReader([]byte(`{"tokn": "abcd"}`))
	_, err := LoadConfig(content)
	if err == nil {
		t.Fatalf("LoadConfig returns empty err")
	}
}
