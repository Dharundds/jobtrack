# jobtrack

Self-hosted job application tracker: sync job-related Gmail, classify messages, and manage your pipeline in a web UI.

**Stack:** Go (Gin, SQLite, Goose) · React (Vite, Tailwind, TanStack Query) · Gmail OAuth + API

## Features

- **Gmail OAuth** — Read-only access (`gmail.readonly`), encrypted refresh tokens in SQLite
- **Sync** — Full sync (90-day window) and incremental sync via Gmail History API; manual **Sync now** from the UI
- **Classification** — Regex rules for applied, rejection, interview, offer, recruiter outreach; filters for LinkedIn noise and weak rejections
- **Applications** — List, kanban board (drag to update status), detail view with timeline and messages
- **Manual cleanup** — Split a message into its own application, remove unrelated mail, delete a listing (ignored messages are not re-grouped on sync)
- **Settings** — Gmail search query, sync interval (stored in DB), classification rules CRUD
- **Stats** — Totals, response rate, counts by status

## Prerequisites

- Go **1.27+**
- Node.js **20+** (for the web UI in dev)
- A [Google Cloud](https://console.cloud.google.com/) project with **Gmail API** enabled and an **OAuth 2.0 Web client**
- Authorized redirect URI: `http://localhost:8080/auth/callback` (dev)

## Quick start

Run the backend and frontend in **two terminals**. Steps below are the same on all platforms; only shell syntax differs.

### 1. Clone and configure

<details>
<summary><strong>Windows</strong> (PowerShell)</summary>

```powershell
cd path\to\jobtrack
Copy-Item .env.example .env
notepad .env
```

Generate `TOKEN_ENC_KEY` and `SESSION_SECRET` (run twice, paste into `.env`):

```powershell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))
```

Or, if [OpenSSL](https://slproweb.com/products/Win32OpenSSL.html) is on your `PATH`:

```powershell
openssl rand -base64 32
```

</details>

<details>
<summary><strong>Linux &amp; macOS</strong></summary>

```bash
cd path/to/jobtrack
cp .env.example .env
${EDITOR:-nano} .env
```

Generate secrets (run twice):

```bash
openssl rand -base64 32
```

</details>

Fill in Google OAuth values in `.env` (see [Environment variables](#environment-variables)).

**Never commit `.env` or `data/`** — they are in `.gitignore`.

### 2. Backend

From the repository root:

| Platform | Command |
|----------|---------|
| **Windows** (PowerShell) | `go run ./cmd/jobtrack` |
| **Linux / macOS** | `go run ./cmd/jobtrack` |

Migrations run on startup. Database: `./data/jobtrack.db` (created automatically).

Optional build:

```powershell
# Windows (PowerShell)
go build -o bin/jobtrack.exe ./cmd/jobtrack
.\bin\jobtrack.exe
```

```bash
# Linux / macOS
go build -o bin/jobtrack ./cmd/jobtrack
./bin/jobtrack
```

### 3. Frontend (development)

Second terminal:

| Platform | Commands |
|----------|----------|
| **Windows** (PowerShell) | `cd web; npm install; npm run dev` |
| **Linux / macOS** | `cd web && npm install && npm run dev` |

Open [http://localhost:5173](http://localhost:5173). Vite proxies `/api` and `/auth` to **http://localhost:8080**.

Production UI build:

```powershell
# Windows (PowerShell)
cd web; npm run build
```

```bash
# Linux / macOS
cd web && npm run build
```

Output: `web/dist/`.

### 4. Connect Gmail

1. Open **Connect** / sign in with Google.
2. After redirect, click **Sync now** in the top bar.
3. Review **Board** or **Applications**.

### Troubleshooting by platform

| Issue | Windows (PowerShell) | Linux / macOS |
|-------|---------|----------------|
| Port 8080 in use | `Get-NetTCPConnection -LocalPort 8080` then stop the process in Task Manager, or change `PORT` in `.env` | `lsof -i :8080` then `kill <pid>`, or change `PORT` |
| `make` not found | Use `go run ./cmd/jobtrack` directly, or install Make via [Git for Windows](https://git-scm.com/download/win) / Chocolatey | Install `make` (`apt install make`, `brew install make`) or use `go` commands below |
| Go module cache errors | `go clean -modcache` then `go mod download` | Same |

Script on Windows: `scripts/fix-go-cache.ps1` if present in the repo.

## Environment variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GOOGLE_CLIENT_ID` | Yes | OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | Yes | OAuth client secret |
| `GOOGLE_REDIRECT_URL` | Yes | Callback URL (default `http://localhost:8080/auth/callback`) |
| `TOKEN_ENC_KEY` | Yes | Base64-encoded 32-byte AES key for refresh tokens |
| `SESSION_SECRET` | Yes | Base64 session signing secret (≥32 bytes decoded) |
| `PORT` | No | HTTP port (default `8080`) |
| `FRONTEND_URL` | No | Post-login redirect in dev (default `http://localhost:5173`) |
| `SYNC_INTERVAL_MINUTES` | No | Background sync interval at process start (default from env) |
| `APP_ENV` | No | `dev` or `production` |
| `LOG_LEVEL` / `LOG_FORMAT` | No | Zerolog level and `console` / `json` |
| `GMAIL_UNITS_PER_MIN` | No | Gmail quota pacing (default `2500`) |
| `GMAIL_CONCURRENCY` | No | Parallel message fetches (default `3`) |

Gmail search scope is configurable in **Settings** (`gmail_query`). If empty, a built-in default query is used (see `internal/config/gmail_query.go`).

## Makefile (Linux & macOS)

On **Linux** and **macOS**, from the repo root:

```bash
make dev      # go run ./cmd/jobtrack
make build    # bin/jobtrack
make test     # go test ./...
make vet      # go vet ./...
make migrate  # goose up (requires goose CLI + existing DB)
```

**Windows:** `make` is optional. Use the equivalents:

| Make target | Windows (PowerShell) |
|-------------|-------------------|
| `make dev` | `go run ./cmd/jobtrack` |
| `make build` | `go build -o bin/jobtrack.exe ./cmd/jobtrack` |
| `make test` | `go test ./...` |
| `make vet` | `go vet ./...` |

## Project layout

```text
cmd/jobtrack/          Entrypoint
internal/
  api/                 HTTP routes, handlers, middleware
  auth/                Google OAuth
  classify/            Rules, grouping, post-sync processor
  config/              Env loading
  crypto/              Token encryption
  gmail/               Client, MIME, retries, quota
  store/               SQLite + queries
  sync/                Background sync worker
migrations/            Goose SQL migrations
web/                   React SPA (Vite)
data/                  Local SQLite (gitignored)
```

## API (authenticated)

Session cookie after OAuth. Main routes:

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/health` | Health check |
| GET | `/api/auth/status` | Connection status |
| POST | `/api/sync` | Trigger sync |
| GET | `/api/sync/status` | Sync progress |
| GET/PATCH | `/api/applications` | List / filter |
| GET/PATCH/DELETE | `/api/applications/:id` | Detail / update / delete listing |
| POST | `/api/applications/merge` | Merge duplicates |
| POST | `/api/applications/:id/messages/:mid/split` | New application from one message |
| DELETE | `/api/applications/:id/messages/:mid` | Remove message from tracking |
| GET/POST/PATCH/DELETE | `/api/rules` | Classification rules |
| GET/PATCH | `/api/settings` | Gmail query, sync interval |
| GET | `/api/stats` | Dashboard metrics |

Auth: `GET /auth/login`, `GET /auth/callback`, `POST /auth/logout`.

## Production notes

- Set `APP_ENV=production`, `LOG_FORMAT=json`, and a strong `FRONTEND_URL` if the UI is served separately.
- Embed or reverse-proxy the built SPA (`web/dist`) in front of the Go server when not using Vite dev.
- Rotate OAuth client secret and encryption keys if `.env` or the DB is ever exposed.
- The SQLite DB contains encrypted OAuth tokens and email excerpts — treat `data/` as sensitive.

## Development

**Tests and frontend build** — same Go/npm commands on every OS:

| Task | Windows (PowerShell) | Linux / macOS |
|------|----------------------|---------------|
| Go tests | `go test ./...` | `go test ./...` |
| Go vet | `go vet ./...` | `go vet ./...` |
| Web build | `cd web; npm run build` | `cd web && npm run build` |
| Web lint (if configured) | `cd web; npm run lint` | `cd web && npm run lint` |

Install tooling:

| Tool | Windows (PowerShell) | Linux (Debian/Ubuntu) | macOS |
|------|---------|------------------------|-------|
| Go | [go.dev/dl](https://go.dev/dl/) | `sudo apt install golang` or official tarball | `brew install go` |
| Node.js | [nodejs.org](https://nodejs.org/) or `winget install OpenJS.NodeJS.LTS` | [nodejs.org](https://nodejs.org/) or `fnm` / `nvm` | `brew install node` |

## License

No license file yet — add one before public distribution if you intend to open-source the project.
