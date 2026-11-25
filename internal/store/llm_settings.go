package store

import (
	"context"
	"strconv"
	"strings"

	"github.com/dharundds/jobtrack/internal/crypto"
)

const (
	SettingLLMBaseURL     = "llm_base_url"
	SettingLLMModel       = "llm_model"
	SettingLLMAPIKeyEnc   = "llm_api_key_enc"
	SettingLLMConcurrency = "llm_concurrency"
)

// LLMConfig holds OpenAI-compatible LLM endpoint settings (API key decrypted in memory).
type LLMConfig struct {
	BaseURL string
	Model   string
	APIKey  string
}

func (c LLMConfig) Configured() bool {
	return c.BaseURL != "" && c.Model != ""
}

func (s *Store) LLMConfig(ctx context.Context, tokenKey []byte) (LLMConfig, error) {
	base, err := s.GetSetting(ctx, SettingLLMBaseURL)
	if err != nil {
		return LLMConfig{}, err
	}
	model, err := s.GetSetting(ctx, SettingLLMModel)
	if err != nil {
		return LLMConfig{}, err
	}
	enc, err := s.GetSetting(ctx, SettingLLMAPIKeyEnc)
	if err != nil {
		return LLMConfig{}, err
	}
	var apiKey string
	if enc != "" {
		apiKey, err = crypto.DecryptToken(tokenKey, enc)
		if err != nil {
			return LLMConfig{}, err
		}
	}
	return LLMConfig{BaseURL: base, Model: model, APIKey: apiKey}, nil
}

func (s *Store) SetLLMConfig(ctx context.Context, tokenKey []byte, baseURL, model, apiKeyPlain string, updateAPIKey bool) error {
	if err := s.SetSetting(ctx, SettingLLMBaseURL, baseURL); err != nil {
		return err
	}
	if err := s.SetSetting(ctx, SettingLLMModel, model); err != nil {
		return err
	}
	if !updateAPIKey {
		return nil
	}
	if apiKeyPlain == "" {
		return s.SetSetting(ctx, SettingLLMAPIKeyEnc, "")
	}
	enc, err := crypto.EncryptToken(tokenKey, apiKeyPlain)
	if err != nil {
		return err
	}
	return s.SetSetting(ctx, SettingLLMAPIKeyEnc, enc)
}

// LLMConcurrency returns parallel LLM request limit (default 2, clamped 1–8).
func (s *Store) LLMConcurrency(ctx context.Context, defaultVal int) (int, error) {
	if defaultVal < 1 {
		defaultVal = 2
	}
	raw, err := s.GetSetting(ctx, SettingLLMConcurrency)
	if err != nil {
		return defaultVal, err
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultVal, nil
	}
	n, convErr := strconv.Atoi(raw)
	if convErr != nil || n < 1 {
		return defaultVal, nil
	}
	if n > 8 {
		return 8, nil
	}
	return n, nil
}

func (s *Store) LLMAPIKeyConfigured(ctx context.Context) (bool, error) {
	enc, err := s.GetSetting(ctx, SettingLLMAPIKeyEnc)
	if err != nil {
		return false, err
	}
	return enc != "", nil
}