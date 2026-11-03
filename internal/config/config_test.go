package config

import (
	"encoding/base64"
	"testing"
)

func TestLoad_MissingRequired(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "")
	t.Setenv("GOOGLE_CLIENT_SECRET", "")
	t.Setenv("TOKEN_ENC_KEY", "")
	t.Setenv("SESSION_SECRET", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing vars")
	}
}

func TestLoad_Valid(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	secret := make([]byte, 32)
	for i := range secret {
		secret[i] = byte(i + 1)
	}

	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	t.Setenv("TOKEN_ENC_KEY", base64.StdEncoding.EncodeToString(key))
	t.Setenv("SESSION_SECRET", base64.StdEncoding.EncodeToString(secret))

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "8080" {
		t.Fatalf("port: got %q", cfg.Port)
	}
	if cfg.ListenAddr() != "127.0.0.1:8080" {
		t.Fatalf("listen: got %q", cfg.ListenAddr())
	}
}

func TestListenAddr_BindAll(t *testing.T) {
	key := make([]byte, 32)
	secret := make([]byte, 32)
	t.Setenv("BIND_ALL", "true")
	t.Setenv("GOOGLE_CLIENT_ID", "x")
	t.Setenv("GOOGLE_CLIENT_SECRET", "y")
	t.Setenv("TOKEN_ENC_KEY", base64.StdEncoding.EncodeToString(key))
	t.Setenv("SESSION_SECRET", base64.StdEncoding.EncodeToString(secret))

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ListenAddr() != ":8080" {
		t.Fatalf("bind all: got %q", cfg.ListenAddr())
	}
}
