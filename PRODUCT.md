# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

**Primary:** Job seekers who run **jobtrack** on their own machine (self-hosted). They already manage applications largely through email (ATS notifications, recruiters, LinkedIn, job boards) and want a single pipeline view without retyping every application into a SaaS tracker.

**Situation:** Active or recent search; Gmail is the main channel for application updates; tolerance for connecting read-only OAuth and tuning sync/query/rules when mail is missed or misclassified.

## Product Purpose

**jobtrack** syncs job-related Gmail, classifies messages into application statuses, and presents applications in a web UI (dashboard, list, kanban board, detail timeline). Success means: mail stays the source of truth with minimal manual data entry; the pipeline reflects reality after sync; users can fix wrong groupings (split, remove, delete) without fighting re-sync; stats and settings support ongoing search hygiene.

## Positioning

A neighboring product cannot truthfully copy all three at once without the same architecture:

1. **Inbox as source of truth** — Gmail sync + rule-based classification instead of manual CRM entry.
2. **Pipeline control** — Board, statuses, detail views, and manual cleanup when automation errs.
3. **Privacy / self-hosting** — Data and tokens on the user’s machine; read-only Gmail scope; no multi-tenant SaaS assumption.

Future messaging and design should **balance** these three; none is the sole hero claim.

## Operating Context

- **Runtime:** Go API on localhost (default `:8080`); React (Vite) UI in dev on `:5173` with proxy to API; SQLite database under `./data/`.
- **Onboarding:** Google Cloud project, Gmail API, OAuth web client; user connects Gmail once via `/connect`; optional tuning of Gmail query and classification rules in Settings.
- **Sync ritual:** Full sync window on first/expired history; incremental sync via Gmail History API; manual **Sync now** from the UI.
- **Dev:** Two terminals (backend + `web` dev server); secrets in `.env` (never committed).

## Capabilities and Constraints

**Capabilities (confirmed in product):**

- Gmail OAuth (read-only `gmail.readonly`); encrypted refresh tokens in SQLite.
- Configurable Gmail search query (seed + Settings); default keywords/ATS senders and subject-based screening terms.
- Classification via regex rules (applied, rejection, interview, offer, outreach, assessments/screening); noise filters (e.g. LinkedIn social, weak rejections).
- Applications: list, kanban (status updates), detail with messages/timeline.
- Manual cleanup: split message to new application, remove unrelated message, delete listing; ignored messages not re-grouped on sync.
- Settings: query, sync interval (DB), rules CRUD.
- Dashboard stats: totals, response rate, counts by status.

**Constraints:**

- Single connected Gmail account model (personal/self-hosted instance, not multi-user SaaS).
- Classification quality depends on rules and extract heuristics; user correction is expected.
- Production UI may be served separately from API (`web/dist` build); embed-in-binary not assumed shipped.

**Terminology:** *Application* = one company/role pipeline; *message* = synced Gmail message linked (or splittable) to applications; *status* = pipeline stage driven by classification and board updates.

## Brand Commitments

- Product name: **jobtrack** (repository folder may differ; module `github.com/dharundds/jobtrack`).
- Voice: practical, trustworthy, no hype; respect that mail content is sensitive.
- No fabricated testimonials, customers, or benchmarks in product copy unless real evidence is added later.

## Evidence on Hand

- **README.md** — setup, features, env vars, platform commands.
- **No** marketing imagery, case studies, or customer quotes checked in; future work must not invent them.
- Real usage data lives only in the user’s local DB and Gmail (not in repo).

## Product Principles

1. **Mail first** — Prefer sync and classification over empty forms; manual edits fix exceptions, not the default workflow.
2. **User owns the data** — Self-hosted storage, explicit OAuth, readable settings for query and rules.
3. **Correctable automation** — Split, remove, delete, and ignore are first-class when rules misfire.
4. **Honest pipeline** — Status and stats should reflect what Gmail supports; avoid implying guarantees (e.g. perfect ATS coverage).
5. **Operate mode UI** — Dense task UI (board, list, settings) prioritizes scanability and consistency over marketing flair.

## Accessibility & Inclusion

Target **WCAG 2.1 Level AA** where practical for the web UI (contrast, focus, labels, keyboard paths on primary flows). Native mobile apps are out of scope unless platform changes.
