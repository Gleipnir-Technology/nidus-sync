# nidus-sync — Cleanup Tasks

This file lists code, files, and patterns that are remnants of older architectural approaches. These should be removed to reduce complexity, maintenance burden, and confusion.

---

## 1. Bob → Jet Migration (Incomplete)

**Status:** Bob is still the primary ORM. Jet was introduced May 2026 but only covers 3 schemas partially.

### 1a. Port remaining schemas from Bob to Jet

Jet-based queries exist for:
- `db/query/public/` — address, communication, communication_log_entry, compliance_report_request, feature, feature_pool, job, lead, signal, site
- `db/query/publicreport/` — compliance, image, image_exif, nuisance, report, report_image, report_log, water
- `db/query/arcgis/` — account, oauth, service_feature, service_map, user, user_privileges

Still using Bob directly (not yet ported to Jet queries):
- `platform/report/notification.go` (13 bob references)
- `platform/background/background.go` (8)
- `platform/arcgis.go` (8)
- `platform/text/send.go` (7)
- `platform/report/some_report.go` (6)
- `platform/site.go` (5)
- `platform/csv/flyover.go` (7)
- `platform/csv/pool.go` (5)
- `platform/csv/csv.go` (4)
- `platform/text/report.go` (4)
- `platform/text/phone_number.go` (3)
- `platform/publicreport/log.go` (3)
- `platform/mailer.go` (3)
- `platform/email/template.go` (2)
- `db/connection.go` (4 — bob.Tx types)
- `db/prepared.go` (2)
- `resource/review_task.go` (2)
- `rmo/status.go` (2)
- `rmo/report.go` (1)
- `rmo/mailer.go` (1)
- Plus many api/* files

### 1b. Remove Bob-generated models after migration

Once all queries are ported to Jet, delete the 103 `.bob.go` files in `db/models/`:
```
db/models/*.bob.go
```

### 1c. Remove Bob-specific helper files

These are Bob-specific and can be removed once Bob is fully replaced:
- `db/dberrors/` — Bob error types (still referenced)
- `db/dbinfo/` — Bob type info (still referenced)
- `db/models/bob_loaders.bob.go`
- `db/models/bob_where.bob.go`

### 1d. Remove Bob from go.mod and dependencies

After all Bob code is gone:
- Remove `github.com/Gleipnir-Technology/bob` from `go.mod`
- Run `go mod tidy`

### 1e. Remove Bob codegen scripts

- `db/bobgen.sh`
- `db/bobgen.yaml`

### 1f. Regenerate Jet output

The `db/jet/main.go` generator outputs to `db/gen/` but no output is currently checked in. Run the generator and ensure generated code is usable:
```bash
cd db/jet && go run .
```

---

## 2. Go HTML Templates → Vue SPA (Mostly Complete)

**Status:** Nearly all Go template routes are commented out in `sync/routes.go` and `rmo/routes.go`. Both hosts serve Vue SPAs via `static.SinglePageApp()`. Some Go template routes remain active.

### 2a. Remaining active Go template routes (sync)

These routes in `sync/routes.go` still render Go templates:
- `/oauth/arcgis/begin` → `getArcgisOauthBegin` (redirect, no template but in Go)
- `/oauth/arcgis/callback` → `getArcgisOauthCallback`
- `/mailer/pool/random` → `getMailerPoolRandom`
- `/mailer/mode-1` → `getMailer1` (generates PDF)
- `/mailer/mode-2` → `getMailer2` (generates PDF)
- `/mailer/mode-3/{code}` → `getMailer3` (generates PDF)
- `/mailer/mode-1/preview` → `getMailer1Preview`
- `/mailer/mode-2/preview` → `getMailer2Preview`
- `/mailer/mode-3/{code}/preview` → `getMailer3Preview`
- `/privacy` → `getPrivacy`

The mailer routes use `platform/pdf` which in turn uses headless Chrome (`chromedp`) to render HTML to PDF. This is legitimate server-side functionality, not just a template remnant. However, the PDF templates themselves may be candidates for migration to the Vue ecosystem.

### 2b. Remove all commented-out routes

Both `sync/routes.go` and `rmo/routes.go` have large blocks of commented-out route registrations. Remove these once migration is confirmed complete.

### 2c. Remove unused Go template files

Once all routes are ported or confirmed dead, remove the entire `html/template/` directory. The `html/` package (`html/embed.go`, `html/filesystem.go`, `html/func.go`, etc.) should also be removed once nothing references it.

### 2d. Reduce the html/ package surface

**Note:** The `html/` package is still actively imported by 40+ Go files. It provides:
- Template rendering (`html/embed.go`, `html/filesystem.go`) — mostly for mailer PDFs and privacy page
- `html.ContentConfig` — used extensively in sync/routes (mailer previews, admin pages)
- `html.MakeGet`, `html.MakePost` — HTTP handler wrappers (used by active `sync/` routes)
- `html.RespondError` — HTTP error responses
- Form parsing, image upload handling, URL building

**Short-term:** Remove the template rendering portion once mailer PDFs and privacy page are migrated.
**Long-term:** The full `html/` package can be removed only after all server-rendered pages are gone and handler wrappers are replaced with the `resource/` pattern.

---

## 3. esbuild (`build.js`) — Removed ✅

*(Completed 2026-05-09: `build.js` removed and `pkgs.esbuild` dropped from flake.nix devShell — Vite is the build tool)*

---

## 4. Legacy Static JavaScript Files

**Status:** `static/js/` contains 20 plain JavaScript files written as custom HTML elements and standalone scripts for the Go template era. These are referenced by old Go HTML templates but most of those templates are now unused.

### 4a. Files in static/js/

```
address-display.js
address-or-report-suggestion.js
address-suggestion.js
events.js
geocode.js
location.js
map-admin.js
map-aggregate.js
map-arcgis-tile.js
map-cell.js
map-locator.js
map-locator-ro.js
map-multipoint.js
map-proxied-arcgis-tile.js
map-routing.js
map-service-area.js
photo-upload.js
table-report.js
table-site.js
time-relative.js
user-selector.js
```

### 4b. Determine which are still used

The remaining active Go templates (mailer, oauth, privacy) may reference some of these. Check each active template for `<script src="/static/js/...">` references. Templates that are confirmed unused:
- All templates in `html/template/sync/` (dashboard, cell, communication-root, district, intelligence, layout, operations-root, planning-root, radar, review, sudo, upload-*) — these are replaced by Vue SPAs
- Most templates in `html/template/rmo/` — RMO routes are all commented out

### 4c. Migrate any still-needed functionality

The map-locator, address-suggestion, and photo-upload functionality has Vue equivalents in `ts/components/`. The remaining custom element patterns should be fully replaced by Vue components.

---

## 5. TomTom Integration — Removed ✅

*(Completed 2026-05-09: `tomtom/` directory removed — zero imports outside itself, Stadia Maps is now the geocoding/tile provider)*

---

## 6. Postgrid — Alternate Mail Provider

**Status:** `postgrid/` contains a single CLI tool (`cmd/send-pdf`) and a `postgrid` Go package reference in `main.go`. Lob is now the mail provider, with its own integration in `lob/`.

### 6a. Investigate and remove if unused

- Check if Postgrid is actually being used in production vs Lob
- If Lob is the chosen provider, remove `postgrid/` entirely
- Remove any Postgrid configuration references

---

## 7. Duplicate Architecture: `api/` vs `resource/`

**Status:** The `api/` package contains both route registration (`api/routes.go`) and handler functions (`api/signin.go`, `api/publicreport.go`, `api/compliance.go`, etc.). The `resource/` package provides typed resource handlers that expose `List`, `Get`, `Create`, etc. Some functionality exists in both layers.

### 7a. Consolidate handler functions

Functions in `api/` that directly handle business logic should be moved to `resource/`:
- `api/signin.go` — `postSignin`, `postSignout`, `postSignup`
- `api/compliance.go` — various compliance handlers
- `api/publicreport.go` — `postPublicreportInvalid`, `postPublicreportSignal`, `postPublicreportMessage`
- `api/sudo.go` — `postSudoEmail`, `postSudoSMS`, `postSudoSSE`
- `api/configuration.go` — `postConfigurationIntegrationArcgis`
- `api/review.go` — `postReviewPool`
- `api/twilio.go`, `api/voipms.go` — webhook handlers
- `api/audio.go`, `api/image.go` — media upload handlers
- `api/tile.go`, `api/debug.go` — utilities

### 7b. Standardize on resource pattern

Either move everything to `resource/` or keep both but clearly define responsibilities:
- `resource/` — domain resource CRUD + URI generation
- `api/` — route registration + HTTP concerns only

Currently the split is unclear and some `api/` files do substantial business logic.

---

## 8. `arcgis-go` Submodule — Not Checked Out

**Status:** The `arcgis-go` submodule (referenced in `.gitmodules`) is not checked out (empty directory). The external `github.com/Gleipnir-Technology/arcgis-go` package is used via `go.mod` instead.

### 8a. Remove submodule

```bash
git submodule deinit arcgis-go
git rm arcgis-go
```

Verify that all code references use the external package, not a local path.

---

## 9. `go-geojson2h3` Local Copy

**Status:** `go-geojson2h3/` is also a submodule. The external package `github.com/Gleipnir-Technology/go-geojson2h3/v2` is imported in `go.mod`. Only `h3utils/h3.go` references it.

### 9a. Consolidate

- If the local copy isn't needed, remove the submodule
- If local modifications exist, merge upstream or maintain intentionally with documentation

---

## 10. Old Generated Files & Artifacts

### 10a. `query.go` at project root

Contains a commented-out Bob query interface and an unused `QueryWriter` interface. The `insertQueryToString` function is entirely commented out. Either repurpose or remove.

### 10b. `db/sql/` directory

Contains `.bob.go` and `.bob.sql` files — these are Bob-style named queries. Once Bob is removed, these can be cleaned up or migrated to Jet equivalents.

### 10c. `static/gen/main.js`

A leftover built artifact. The new build output goes to `static/gen/sync/` and `static/gen/rmo/` via Vite. Ensure `static/gen/` is in `.gitignore` and the stale `main.js` is removed.

### 10d. `static/css/placeholder`

Empty placeholder file. Remove.

---

## 11. Nix devShell Cleanup

**Status:** `flake.nix` devShell includes several tools from older workflows:

### 11a. Potentially unnecessary devShell packages

- `pkgs.esbuild` — replaced by Vite (keep only if `build.js` is retained)
- `pkgs.dart-sass` — Vue/Vite uses the `sass` npm package; check if Go code invokes dart-sass directly
- `pkgs.autoprefixer` — may not be needed with Vite's built-in PostCSS

---

## 12. Start Scripts — Consolidate

**Status:** Four start scripts exist:

| Script | Purpose |
|--------|---------|
| `start-air.sh` | Development with air (live reload) |
| `start-flogo.sh` | Unknown (references `flogo`) |
| `start-nidus-sync.sh` | Production-like direct run |
| `start-nix-built.sh` | Run Nix-built output |

`start-flogo.sh` may be a remnant. Investigate and remove if unused.

---

## Priority Summary

1. **High impact, low effort:**
   - ~~Remove `tomtom/` (unused, no imports)~~ ✅
   - ~~Remove `build.js` (dead, replaced by Vite)~~ ✅
   - Remove commented-out routes in `sync/routes.go` and `rmo/routes.go`
   - Remove `query.go` commented-out code
   - Remove `static/gen/main.js` stale artifact
   - Remove `static/css/placeholder`

2. **Medium impact, medium effort:**
   - Remove unused Go HTML templates (confirm which are still active first)
   - Remove unused `static/js/` files (verify against active templates)
   - Remove `arcgis-go` submodule
   - Clean up Nix devShell

3. **High impact, high effort:**
   - Complete Bob → Jet migration across all schemas
   - Remove Bob-generated models, helpers, scripts
   - Remove Bob from go.mod
   - Consolidate `api/` and `resource/` handler patterns
   - Remove `html/` package (after all Go templates are gone)
