# AGENTS.md — Nidus Sync Codebase Guide

This file captures conventions, patterns, and gotchas for anyone working on this codebase. It was produced during a lint cleanup pass (May 2026) to document lessons learned.

## Project Overview

**Module:** `github.com/Gleipnir-Technology/nidus-sync`
**Language:** Go 1.24
**Build:** Nix (`flake.nix`) + standard Go toolchain
**ORM:** Bob (legacy) + Jet (new, partial migration)
**Frontend:** Vue SPA (Vite) replacing Go HTML templates

The app serves two hosts from a single binary:
- **Sync** (`sync/`) — internal dashboard for mosquito control districts
- **RMO** (`rmo/`) — public-facing "Report Mosquitoes Online" site

Both are migrating from Go `html/template` rendered pages to Vue SPAs served by `static.SinglePageApp()`.

## Build & Lint Commands

```bash
# Build everything
go build ./...

# Run linter
golangci-lint run

# Build a specific package
go build ./api/
go build ./platform/
```

## Lint Helpers (`lint/error.go`)

The `lint/` package provides helpers for common error-handling patterns. **Always use these instead of bare calls** to avoid errcheck lint failures:

| Helper | Use for | Example |
|--------|---------|---------|
| `lint.Fprintf(w, fmt, args...)` | `fmt.Fprintf` to writers where errors are non-critical | `lint.Fprintf(w, "ok")` |
| `lint.Fprint(w, args...)` | `fmt.Fprint` to writers | `lint.Fprint(w, "User-agent: *\n")` |
| `lint.Write(w, p []byte)` | `w.Write(p)` — HTTP response bodies | `lint.Write(w, body)` |
| `lint.LogOnErr(f, msg)` | Deferred `Close()` calls | `defer lint.LogOnErr(file.Close, "close file")` |
| `lint.LogOnErrCtx(f, ctx, msg)` | `txn.Commit(ctx)` or other ctx funcs | `lint.LogOnErrCtx(txn.Commit, ctx, "commit")` |
| `lint.LogOnErrRollback(f, ctx, msg)` | Deferred `txn.Rollback(ctx)` | `defer lint.LogOnErrRollback(txn.Rollback, ctx, "rollback")` |

**Key rule:** `LogOnErrRollback` silently ignores `"sql: transaction has already been committed or rolled back"` errors, which occur when a deferred rollback fires after a successful commit. Always use it for deferred rollbacks.

### For DB transactions, use this pattern:

```go
txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
if err != nil {
    return fmt.Errorf("begin: %w", err)
}
defer lint.LogOnErrRollback(txn.Rollback, ctx, "rollback")

// ... do work ...

if err := txn.Commit(ctx); err != nil {
    return fmt.Errorf("commit: %w", err)
}
return nil
```

### For HTTP handlers that render HTML:

```go
if err := renderShim(w, r, errRender(err)); err != nil {
    http.Error(w, fmt.Sprintf("render shim: %v", err), http.StatusInternalServerError)
}
```

## Architecture Notes

### Two hosts, one binary

`main.go` creates two `gorilla/mux` routers and supports three modes via CLI flags:
- `-sync` — serve the Sync dashboard
- `-report` — serve the RMO public site
- `-all` — serve both (default)

Each host has its own route registration in `sync/routes.go` and `rmo/routes.go`.

### RMO package — all handlers are dead

**All** route registrations in `rmo/routes.go` are commented out. The file now only serves the Vue SPA via `static.SinglePageApp("static/gen/rmo")`. During cleanup (May 2026), all handler files were deleted:

```
rmo/compliance.go  rmo/error.go     rmo/nuisance.go  rmo/report.go
rmo/district.go    rmo/image.go     rmo/quick.go     rmo/root.go
rmo/email.go       rmo/mailer.go    rmo/notification.go  rmo/scss.go
rmo/mock.go        rmo/status.go    rmo/water.go
```

Only `rmo/routes.go` remains. **Do not add new Go template handlers here** — the RMO host is pure Vue SPA now.

### Sync package — partially live

Many route registrations in `sync/routes.go` are active. Files deleted during cleanup were those with zero active registrations:

```
sync/admin.go       sync/download.go      sync/operations.go   sync/pool.go
sync/cell.go        sync/intelligence.go  sync/parcel.go       sync/radar.go
sync/communication.go  sync/messages.go   sync/planning.go     sync/review.go
sync/dash.go        sync/mock.go          sync/notification.go sync/service-request.go
sync/signin.go      sync/sms.go           sync/text.go         sync/tile.go
```

### api/ vs resource/ — two handler layers

The codebase has two HTTP handler patterns:

1. **`api/`** — route registration (`api/routes.go`) + `http.HandlerFunc` handlers. Handles signin, webhooks (Twilio, VoIP.ms), media uploads, configuration POSTs.

2. **`resource/`** — typed resource handlers with `List`, `Get`, `Create`, etc. methods. Each resource has a struct embedding `*router` for URI generation. This is the newer, preferred pattern.

The split is not clean — some `api/` files contain substantial business logic. New handlers should use the `resource/` pattern.

### DB access — Bob vs Jet

Two ORMs coexist:
- **Bob** (`github.com/Gleipnir-Technology/bob`) — legacy, used by most queries. Models in `db/models/*.bob.go` (103 files).
- **Jet** (`db/jet/`) — new, generated queries in `db/query/public/`, `db/query/publicreport/`, `db/query/arcgis/`. Only 3 schemas partially ported.

The `db.PGInstance` singleton holds both `BobDB` and `PGXPool`. Jet uses PGXPool directly; Bob uses BobDB.

### db/prepared.go & db/fieldseeker.go

`db/prepared.go` contains utility functions (`pointOrNull`, `lineOrNull`, `queryStoredProcedure`, etc.) that are **only** called from `db/fieldseeker.go`. That file is **entirely commented out** (`/* ... */`). The 9 unused-prepared-funcs lint warnings are expected — do not delete them unless you're also deleting or uncommenting fieldseeker.go.

## Lint Cleanup Context (May 2026)

### What was fixed

- **errcheck (36→0):** All unchecked error returns eliminated using `lint/` helpers or explicit checks.
- **unused (50→9):** ~60 functions/types deleted across ~30 files. Remaining 9 are in `db/prepared.go` (see above).

### golangci-lint reporting cap

golangci-lint caps unused reports at **50 items**. During cleanup, each batch of deletions exposed previously hidden items. If you see 50 unused items, there are almost certainly more hidden behind the cap. Delete the visible ones, re-run lint, and repeat.

### What was NOT fixed (remaining lint categories)

- **govet (26):** printf format mismatches, copylocks, lostcancel — some are real bugs
- **ineffassign (9):** dead assignments that may indicate logic errors
- **staticcheck (29):** deprecated `io/ioutil`, redundant returns, error string conventions, comparison always-true bugs

### Deleted by file count

| Directory | Files deleted | Reason |
|-----------|--------------|--------|
| `rmo/` | 15 | All handlers unused — routes commented out |
| `sync/` | 17 | Unregistered handlers |
| `api/` | 2 | `compliance.go`, `debug.go` — unused handlers |
| `platform/` | 2 | `text/db.go`, `dashboard.go`, `publicreport/address.go` |
| Other | 1 | `tomtom/` (prior cleanup) |

## Commit Conventions

Commits during the cleanup followed a consistent pattern:

```
lint: fix errcheck in api/api.go debug log writes
lint: remove unused code from sync/ package
```

Each commit fixes one category of issue in a small set of related files. Build verification (`go build ./...`) was performed before each commit.

## See Also

- `CLEANUP.md` — broader cleanup roadmap (Bob→Jet migration, html/ package removal, etc.)
- `HISTORY.md` — project history and architectural decisions
- `README.md` — administration and build-from-source instructions
