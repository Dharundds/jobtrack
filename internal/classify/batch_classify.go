package classify

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/dharundds/jobtrack/internal/store"
	"github.com/rs/zerolog"
)

func emailInputFromMessage(m store.MessageRow, accountEmail string) EmailInput {
	return EmailInput{
		From:         m.FromAddr,
		To:           m.ToAddr,
		Subject:      m.Subject,
		Body:         m.BodyExcerpt,
		Snippet:      m.Snippet,
		Direction:    m.Direction,
		AccountEmail: accountEmail,
		ReceivedAt:   m.ReceivedAt,
	}
}

// ClassifyMessages classifies messages in parallel when LLM is enabled (bounded by opts.LLMConcurrency).
func ClassifyMessages(
	ctx context.Context,
	log zerolog.Logger,
	msgs []store.MessageRow,
	accountEmail string,
	ruleClassifier *RuleClassifier,
	opts RunOptions,
) map[int64]Result {
	p := &Processor{log: log.With().Str("component", "classify").Logger()}
	out := make(map[int64]Result, len(msgs))
	total := len(msgs)
	if total == 0 {
		return out
	}

	report := func(done, llmN int) {
		if opts.OnProgress == nil {
			return
		}
		opts.OnProgress(Progress{
			Phase:              "classify",
			ClassifyTotal:      total,
			ClassifyDone:       done,
			LLMReviewed:        llmN,
			ApplicationsLinked: 0,
		})
	}

	if !opts.LLMAvailable || opts.LLM == nil {
		for i, m := range msgs {
			if ctx.Err() != nil {
				break
			}
			in := emailInputFromMessage(m, accountEmail)
			out[m.ID] = ruleClassifier.Classify(ctx, in)
			report(i+1, 0)
		}
		return out
	}

	conc := opts.LLMConcurrency
	if conc < 1 {
		conc = 2
	}
	if conc > 8 {
		conc = 8
	}

	prefer := NewPreferLLM(true, opts.LLM, ruleClassifier, p.log)
	sem := make(chan struct{}, conc)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var doneCnt atomic.Int32
	var llmCnt atomic.Int32

	for _, m := range msgs {
		m := m
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if ctx.Err() != nil {
				return
			}
			in := emailInputFromMessage(m, accountEmail)
			res := prefer.Classify(ctx, in)
			mu.Lock()
			out[m.ID] = res
			mu.Unlock()

			d := int(doneCnt.Add(1))
			if res.ClassifiedBy == "llm" {
				llmCnt.Add(1)
			}
			report(d, int(llmCnt.Load()))
		}()
	}
	wg.Wait()
	return out
}

func countLLMInMap(m map[int64]Result) int {
	n := 0
	for _, r := range m {
		if r.ClassifiedBy == "llm" {
			n++
		}
	}
	return n
}