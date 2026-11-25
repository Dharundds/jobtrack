package classify

import (
	"context"

	"github.com/rs/zerolog"
)

// PreferLLM uses the LLM when enabled; falls back to rules on error.
type PreferLLM struct {
	LLM     LLMClassifier
	Rules   *RuleClassifier
	Enabled bool
	Log     zerolog.Logger
}

func NewPreferLLM(enabled bool, client LLMClassifier, rules *RuleClassifier, log zerolog.Logger) *PreferLLM {
	return &PreferLLM{
		LLM:     client,
		Rules:   rules,
		Enabled: enabled,
		Log:     log,
	}
}

func (p *PreferLLM) Classify(ctx context.Context, in EmailInput) Result {
	if p.Enabled && p.LLM != nil {
		res, err := p.LLM.Classify(ctx, in)
		if err == nil {
			return res
		}
		p.Log.Warn().Err(err).Msg("llm classify failed, using rules")
	}
	return p.Rules.Classify(ctx, in)
}