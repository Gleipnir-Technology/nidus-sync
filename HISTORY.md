# nidus-sync — Project History

## Overview

nidus-sync is a dual-tenant mosquito abatement platform serving two domains:
- **RMO** (`report.mosquitoes.online`) — Public-facing mosquito/water/nuisance reporting
- **Sync** (`sync.nidus.cloud`) — Administrative dashboard for vector control districts

The project was started in November 2025 and has undergone several major architectural shifts across ~1655 commits spanning 6 months.

---

## Timeline

### Phase 1: Foundation (November 2025)

**Nov 3 – Nov 13: Project bootstrap**
- Initial Go project with Nix build system (`flake.nix`, `default.nix`)
- Basic `net/http` web serving with `gorilla/mux` routing
- Go `html/template` server-side rendering
- Bob ORM integration (`github.com/Gleipnir-Technology/bob`) for PostgreSQL — code-generated models via `bobgen`
- ArcGIS OAuth integration for user authentication
- ArcGIS Fieldseeker data synchronization (treatment areas, inspections, breeding sources, etc.)
- MapBox GL JS integration for heatmap visualization
- Dashboard with login, basic CRUD mocks

**Nov 13 – Nov 24: Logging & DB restructuring**
- Migration from standard `log` to `zerolog` for structured, colorized output
- Database logic moved into a separate `db/` subdirectory
- Clean shutdown logic, token refresh loops

**Key characteristics:** Monolithic Go server, HTML templates, Bob ORM, MapBox maps, ArcGIS OAuth

---

### Phase 2: Fieldseeker & Schema Evolution (December 2025)

**Dec 2 – Dec 24: Fieldseeker schema v2**
- Bob codegen updated to latest version
- Fieldseeker schema captured on OAuth connect and stored locally
- Dynamic SQL functions replacing hardcoded per-table sync logic
- Old Fieldseeker tables removed, v2 generated tables used
- Note/image audio support added
- MMS file downloads from SMS webhooks

**Key characteristics:** Bob-generated fieldseeker models, prepared SQL functions, SMS/MMS debugging

---

### Phase 3: Architecture Maturation (January 2026)

**Jan 2 – Jan 8: Domain split & template system**
- WIP pass-through models concept ("Checkpoint on initial idea for passing through models")
- Massive reorganization: templates split into `rmo/` (public) and `sync/` (admin) subdirectories
- `html/` package created with embedded template loading
- Bob submodule removed, `arcgis-go` became external dependency
- Public report domain support added
- Version bumped 7 times in rapid iteration (v0.0.4 → v0.0.10)

**Jan 8 – Jan 31: Platform Layer emergence**
- "Report platform layer" introduced (`a9b0a55f`) — initial abstraction between HTTP handlers and database
- Address suggestion and map-locator components via custom HTML elements
- SVG auto-transformation into Go templates
- Report submission forms wired up (nuisance, water)
- Email template system

**Key characteristics:** Two-domain architecture (RMO/Sync), `html/` template package, platform layer beginning, custom element web components

---

### Phase 4: Map Migration & Platform Expansion (February 2026)

**Feb 1 – Feb 28: Map provider transition**
- MapBox → MapLibre GL (open-source fork) via `maplibre-gl`
- Stadia Maps integration for tile serving and geocoding (Feb 12-14)
- TomTom routing integration added (Feb 17)
- Bulk geocoding via Stadia
- Parcel image generation debugging

**Platform layer expansion:**
- Emails moved to platform layer
- Phone/SMS support
- OAuth integration settings
- Upload platform functions
- QR code and image tile moved into platform
- Admin map components

**Key characteristics:** MapLibre/Stadia replacing MapBox, TomTom added, platform layer expanding, heavy template iteration

---

### Phase 5: VueJS Revolution (March 2026) — 448 commits

**Mar 5 – Mar 12: Pre-Vue cleanup**
- Stadia Maps client initialization
- Signal database schema added
- Review task/mailer schema rework
- Generated Bob files pruned

**Mar 12: Massive platform layer rework** (`44c4f17f`)
- User/organization handling restructured in platform layer
- Signal creation moved inside platform

**Mar 18 – Mar 22: VueJS Migration** (the biggest architectural shift)
- Mar 18: Auto-generated report IDs
- Mar 21: **VueJS introduced** — begins with TypeScript bundle, then Vue SFC components, vue-router, Bootstrap/SCSS integration
- Mar 21: Dashboard, Intelligence, sidebar all moved to Vue
- Mar 22: **esbuild replaced by Vite** (`47f900ab`) — `vite/` directory with separate configs for `sync` and `rmo` SPAs
- Mar 22: TypeScript checking clean across entire frontend
- Mar 23: Public report card component, auth checks off API client
- Mar 24-31: Communication page ripped into components, impersonation support, users page

**Key characteristics:** VueJS 3 + TypeScript + Vite frontend, Pinia stores, vue-router, SCSS, SPA architecture replacing server-rendered Go templates

---

### Phase 6: Compliance & Communication (April 2026) — 454 commits

**Apr 1 – Apr 9: RMO frontend & resources**
- Resource layer expanded (user, avatar, district, nuisance, water, compliance resources)
- RMO frontend checkpoint — Vue ports of public-facing pages
- TS types migrated into API module
- Old bundle paths removed, old SPA generation removed

**Apr 10 – Apr 17: Compliance workflow**
- Compliance report creation, mailer flow
- Site/pool review tasks
- Stadia Maps cache, direct tile access
- OAuth refresh in frontend
- Image upload components

**Apr 17 – Apr 25: Communication system**
- Background jobs reworked for shorter transactions
- Lob (physical mail) integration — direct API client, address creation, letter events
- QR code generation moved to API
- Compliance report evidence, mailer views
- Vue map system generalized (`cad01e68`)

**Apr 25 – Apr 30: Map & communication polish**
- VueJS reimplementation of address/report suggestion
- Communication workbench with map, list, detail views
- Text message log, email/phone display
- Compliance card detail display
- SSE event system with status vs resource message distinction
- Systemd socket activation for downtime-free deploys
- Sentry error tracking for Vue frontend

**Key characteristics:** Compliance/mailer operational, communication system born, Lob integration, Sentry, generalized Vue map system

---

### Phase 7: Jet Migration & Cleanup (May 2026) — 46 commits so far

**May 1 – May 9: SQL generation transition**
- **Jet (go-jet/jet) introduced** — type-safe SQL builder replacing Bob's query building
- Custom Jet generator created with geometry/Box2D type support (`db/jet/main.go`)
- `publicreport` schema ported to Jet
- `arcgis` schema ported to Jet (compiles, not fully tested per commit message)
- New `communication` table added
- Communication marking workflow (invalid, pending-response, possible-issue, possible-resolved)
- Linting: `golangci-lint` added to lefthook, per-file linting
- Cleanup of legacy generated columns (latitude/longitude), string-based queries
- Centralized error handler for Vue sync app

**Key characteristics:** Bob→Jet transition in progress, communication workflow, code quality improvements

---

## Architectural Patterns (by layer)

### Current architecture stack

```
┌─────────────────────────────────────────────────┐
│  Vue 3 SPA (TypeScript)                          │
│  ts/ — shared components, composables, stores    │
│  vite/sync/ — admin SPA entry                    │
│  vite/rmo/ — public SPA entry                    │
├─────────────────────────────────────────────────┤
│  Go HTTP Server (gorilla/mux)                    │
│  api/routes.go — central route registration      │
│  resource/ — resource handlers (REST patterns)   │
│  sync/ — remaining Go template routes             │
│  rmo/ — remaining Go template routes              │
├─────────────────────────────────────────────────┤
│  platform/ — business logic layer                │
│  (address, compliance, communication, district,  │
│   email, fieldseeker, mailer, publicreport,      │
│   review, signal, text, user, upload, etc.)      │
├─────────────────────────────────────────────────┤
│  db/ — database access                           │
│  db/models/ — Bob-generated models (103 files)   │
│  db/query/ — Jet-based query functions           │
│  db/prepared.go — prepared SQL functions         │
├─────────────────────────────────────────────────┤
│  PostgreSQL                                      │
└─────────────────────────────────────────────────┘
```

### Pattern: Platform Layer
Introduced January 2026, the `platform/` package encapsulates business logic between HTTP handlers and the database. It grew from initial report handling to encompass users, organizations, emails, texts, compliance, communications, signals, geocoding, tiles, uploads, and more.

### Pattern: Resource Layer
Added March–April 2026, `resource/` provides typed REST resource handlers with URI generation (via mux route naming). Resources are instantiated with a `resource.NewRouter()` and expose methods like `List`, `Get`, `Create`, `Update`, `Delete` that return domain types. This replaced ad-hoc handler functions in `api/`.

### Pattern: Dual SPA + API
Since late March 2026, both domains serve Vue SPAs for most routes, with the Go server acting as an API backend. The `static.SinglePageApp()` handler serves the Vite-built output and falls back to `index.html` for client-side routing. Some Go template routes remain for mailer PDF generation, OAuth flows, and previews.
