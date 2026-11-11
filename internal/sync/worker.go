package sync

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dharundds/jobtrack/internal/auth"
	"github.com/dharundds/jobtrack/internal/classify"
	"github.com/dharundds/jobtrack/internal/config"
	gmailclient "github.com/dharundds/jobtrack/internal/gmail"
	"github.com/dharundds/jobtrack/internal/privacy"
	"github.com/dharundds/jobtrack/internal/store"
	"github.com/rs/zerolog"
)

// Status is exposed via GET /api/sync/status.
type Status struct {
	Running      bool       `json:"running"`
	LastSyncedAt *time.Time `json:"last_synced_at,omitempty"`
	LastError    string     `json:"last_error,omitempty"`
	Processed    int        `json:"processed,omitempty"`
}

type Worker struct {
	cfg   *config.Config
	store *store.Store
	auth  *auth.Service
	log   zerolog.Logger

	mu        sync.RWMutex
	running   bool
	lastError string
	processed int

	trigger chan struct{}
}

func NewWorker(cfg *config.Config, st *store.Store, authSvc *auth.Service, log zerolog.Logger) *Worker {
	return &Worker{
		cfg:     cfg,
		store:   st,
		auth:    authSvc,
		log:     log.With().Str("component", "sync").Logger(),
		trigger: make(chan struct{}, 1),
	}
}

func (w *Worker) Start(ctx context.Context) {
	interval := time.Duration(w.cfg.SyncIntervalMin) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	w.log.Info().
		Dur("interval", interval).
		Int("gmail_units_per_min", w.cfg.GmailUnitsPerMin).
		Int("gmail_concurrency", w.cfg.GmailConcurrency).
		Msg("sync worker started")

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("sync worker stopped")
			return
		case <-ticker.C:
			w.runSync(ctx)
		case <-w.trigger:
			w.runSync(ctx)
		}
	}
}

// ResetSyncHistory clears the Gmail history cursor so the next sync uses the search query (full list).
func (w *Worker) ResetSyncHistory(ctx context.Context) error {
	return w.store.ClearSyncHistory(ctx)
}

// Trigger queues a sync (non-blocking). Returns false if a sync is already running.
func (w *Worker) Trigger() bool {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return false
	}
	w.mu.Unlock()
	select {
	case w.trigger <- struct{}{}:
		return true
	default:
		return false
	}
}

// TriggerManual starts a sync or returns current status if already running.
func (w *Worker) TriggerManual(ctx context.Context) (started bool, status Status) {
	w.mu.Lock()
	if w.running {
		status = w.statusLocked(ctx)
		w.mu.Unlock()
		return false, status
	}
	w.mu.Unlock()

	if !w.Trigger() {
		return false, w.StatusSnapshot(ctx)
	}
	return true, Status{}
}

func (w *Worker) StatusSnapshot(ctx context.Context) Status {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.statusLocked(ctx)
}

func (w *Worker) statusLocked(ctx context.Context) Status {
	st := Status{
		Running:   w.running,
		LastError: w.lastError,
		Processed: w.processed,
	}
	acc, err := w.store.GetAccount(ctx)
	if err == nil && acc.LastSyncedAt.Valid {
		t := acc.LastSyncedAt.Time
		st.LastSyncedAt = &t
	}
	return st
}

func (w *Worker) runSync(parent context.Context) {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return
	}
	w.running = true
	w.lastError = ""
	w.processed = 0
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		w.running = false
		w.mu.Unlock()
	}()

	ctx, cancel := context.WithTimeout(parent, 30*time.Minute)
	defer cancel()

	processed, err := w.syncOnce(ctx)
	w.mu.Lock()
	w.processed = processed
	if err != nil {
		w.lastError = syncErrorMessage(err)
	} else {
		w.lastError = ""
	}
	w.mu.Unlock()

	if err != nil {
		w.log.Warn().Err(err).Int("stored", processed).Msg("sync finished with error")
		return
	}
	w.log.Info().Int("stored", processed).Msg("sync completed")
}

func (w *Worker) syncOnce(ctx context.Context) (int, error) {
	acc, err := w.store.GetAccount(ctx)
	if err != nil {
		return 0, fmt.Errorf("no connected account")
	}

	client, err := gmailclient.NewClientFromAccount(
		ctx, w.auth, w.cfg.TokenEncKey, acc, w.cfg.GmailUnitsPerMin, w.log,
	)
	if err != nil {
		return 0, err
	}

	query, err := w.store.GmailQuery(ctx)
	if err != nil {
		return 0, err
	}

	var ids []string
	fullSync := false

	if acc.HistoryID.Valid && acc.HistoryID.String != "" {
		startID, parseErr := strconv.ParseUint(acc.HistoryID.String, 10, 64)
		if parseErr != nil {
			fullSync = true
		} else {
			var expired bool
			ids, expired, err = client.ListHistory(ctx, startID)
			if err != nil {
				return 0, err
			}
			if expired {
				w.log.Info().Msg("history expired, full resync")
				fullSync = true
				ids = nil
			}
		}
	} else {
		fullSync = true
	}

	if fullSync {
		since := time.Now().UTC().AddDate(0, 0, -config.FullSyncDays)
		ids, err = client.ListMessageIDs(ctx, query, since)
		if err != nil {
			return 0, err
		}
		w.log.Info().Int("listed", len(ids)).Msg("full sync message ids")
	} else {
		w.log.Info().Int("listed", len(ids)).Msg("incremental history message ids")
	}

	newIDs, skipped, err := w.store.FilterNewGmailIDs(ctx, ids)
	if err != nil {
		return 0, err
	}
	newIDs = uniqueStrings(newIDs)
	w.log.Info().
		Int("listed", len(ids)).
		Int("skipped_existing", skipped).
		Int("to_fetch", len(newIDs)).
		Msg("sync dedupe")

	stored, fetchErr := w.fetchAndStore(ctx, client, acc.Email, newIDs)
	if fetchErr != nil {
		// Keep stored messages; do not advance historyId so the next run resumes.
		return stored, fetchErr
	}

	historyID, err := client.CurrentHistoryID(ctx)
	if err != nil {
		return stored, err
	}

	processor := classify.NewProcessor(w.store, w.log)
	linked, procErr := processor.Run(ctx, acc.Email)
	if procErr != nil {
		w.log.Warn().Err(procErr).Int("linked", linked).Msg("classify/group failed")
		return stored, procErr
	}
	w.log.Info().Int("messages_linked", linked).Msg("classify and group")

	syncedAt := time.Now().UTC().Format(time.RFC3339)
	if err := w.store.UpdateAccountSync(ctx, fmt.Sprintf("%d", historyID), syncedAt); err != nil {
		return stored, err
	}

	return stored, nil
}

func (w *Worker) fetchAndStore(
	ctx context.Context,
	client *gmailclient.Client,
	accountEmail string,
	ids []string,
) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	concurrency := w.cfg.GmailConcurrency
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	var stored atomic.Int32
	var failOnce sync.Once
	var failErr error

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fail := func(err error) {
		if err == nil || errors.Is(err, context.Canceled) {
			return
		}
		failOnce.Do(func() {
			failErr = err
			w.log.Warn().Err(err).Msg("sync fetch failed")
			cancel()
		})
	}

	for _, id := range ids {
		if ctx.Err() != nil {
			break
		}
		if failErr != nil {
			break
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(gmailID string) {
			defer wg.Done()
			defer func() { <-sem }()

			if ctx.Err() != nil || failErr != nil {
				return
			}

			msg, err := client.GetMessage(ctx, gmailID)
			if err != nil {
				if gmailclient.IsNotFound(err) {
					w.log.Debug().Str("gmail_id", gmailID).Msg("skip missing message")
					return
				}
				if gmailclient.IsRateLimitError(err) {
					fail(fmt.Errorf("gmail rate limit exceeded after retries: %w", err))
					return
				}
				fail(err)
				return
			}

			direction := "in"
			if hasLabel(msg.LabelIDs, "SENT") || addressesMatch(msg.From, accountEmail) {
				direction = "out"
			}

			_, err = w.store.InsertMessage(ctx, store.MessageInsert{
				GmailID:     msg.ID,
				ThreadID:    msg.ThreadID,
				Direction:   direction,
				FromAddr:    msg.From,
				ToAddr:      msg.To,
				Subject:     msg.Subject,
				Snippet:     msg.Snippet,
				BodyExcerpt: msg.BodyExcerpt,
				ReceivedAt:  msg.ReceivedAt,
			})
			if err != nil {
				if store.IsDuplicateMessage(err) {
					stored.Add(1)
					return
				}
				fail(err)
				return
			}

			n := stored.Add(1)
			if n%25 == 0 {
				w.log.Info().Int("stored", int(n)).Msg("sync progress")
			}
			w.log.Debug().
				Str("gmail_id", msg.ID).
				Str("direction", direction).
				Str("from", privacy.MaskEmail(extractEmail(msg.From))).
				Msg("stored message")
		}(id)
	}

	wg.Wait()

	if failErr != nil {
		return int(stored.Load()), failErr
	}
	return int(stored.Load()), nil
}

func syncErrorMessage(err error) string {
	if gmailclient.IsOAuthTokenError(err) {
		return "Gmail authorization expired — disconnect and connect Gmail again"
	}
	return err.Error()
}

func uniqueStrings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func hasLabel(labels []string, want string) bool {
	for _, l := range labels {
		if l == want {
			return true
		}
	}
	return false
}

func addressesMatch(fromHeader, accountEmail string) bool {
	fromHeader = strings.ToLower(strings.TrimSpace(fromHeader))
	accountEmail = strings.ToLower(strings.TrimSpace(accountEmail))
	return strings.Contains(fromHeader, accountEmail)
}

func extractEmail(from string) string {
	if i := strings.Index(from, "<"); i >= 0 {
		if j := strings.Index(from, ">"); j > i {
			return from[i+1 : j]
		}
	}
	return from
}
