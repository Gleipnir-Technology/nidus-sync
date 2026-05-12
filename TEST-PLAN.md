# TEST-PLAN.md — Nidus Sync Database Query Layer Testing

## Status

The project currently has **zero tests**. `testify` is already in `go.mod` as an indirect dependency and needs to be promoted to direct.

## Overview

This plan covers **Tier 2 testing**: integration tests for the `db/query/` packages that execute real SQL against a throwaway Postgres database using the project's own migration system. All tests operate inside transactions that are rolled back, so they leave no trace.

The query layer is the ideal starting point because:

1. Every function is small and focused — a single INSERT, SELECT, or UPDATE.
2. After the signature normalization (CLEANUP.md §13), **every** function will accept a `db.Ex` interface, making them all uniformly testable from a test transaction.
3. These are the foundation that all platform-layer business logic calls. Bugs here cascade upward.

### Prerequisite: Normalize Query Function Signatures

Before writing tests, all query functions must be converted to accept `(ctx context.Context, txn db.Ex, ...)`. This is documented in detail at **[CLEANUP.md §13](CLEANUP.md#13-normalize-query-function-signatures-to-dbex)**. Summary of changes needed:

| Category | Count | What | Test-blocking? |
|----------|-------|------|---------------|
| 13d — Bugfix: txn ignored | 2 funcs | `AddressFromID`, `AddressFromComplianceReportRequestID` call `ExecuteOne` instead of `ExecuteOneTx` | Yes — data isolation broken |
| 13b — `db.Tx` → `db.Ex` | 4 funcs | `CommunicationInsert`, `CommunicationSetStatus`, `CommunicationLogEntryInsert`, `ComplianceFromID` | Partial — works but can't pass mock |
| 13a — Add `txn db.Ex` param | 19 funcs | Functions missing transaction parameter entirely | Yes — can't test in transactions |
| 13c — `bob.Tx` → `db.Ex` | 6 funcs | ArcGIS package functions using Bob transactions | Yes — can't test without Bob |

**Order of operations:** Fix 13d → convert 13b → convert 13a → convert 13c. After all conversions, every function is testable with `dbtest.Txn()`.

---

## Architecture of the Query Layer

### Package structure

```
db/query/
├── public/          ← main "public" schema queries (Jet ORM)
│   ├── address.go
│   ├── communication.go
│   ├── communication_log_entry.go
│   ├── compliance_report_request.go
│   ├── feature.go
│   ├── feature_pool.go
│   ├── job.go
│   ├── lead.go
│   ├── signal.go
│   └── site.go
├── publicreport/    ← "publicreport" schema queries (Jet ORM)
│   ├── compliance.go
│   ├── image.go
│   ├── image_exif.go
│   ├── nuisance.go
│   ├── report.go
│   ├── report_image.go
│   ├── report_log.go
│   └── water.go
└── arcgis/          ← "arcgis" schema queries (Jet ORM)
    ├── account.go
    └── ...
```

### Query function patterns

There are three patterns in the query layer:

| Pattern | Signature | Example |
|---------|-----------|---------|
| **Insert (txn)** | `func XxxInsert(ctx, txn db.Ex, model) (model, error)` | `CommunicationInsert`, `LeadInsert`, `ReportInsert` |
| **Insert (global)** | `func XxxInsert(ctx, model) (model, error)` | (would use `db.PGInstance` directly) |
| **Select (txn)** | `func XxxFromYyy(ctx, txn db.Ex, ...) (model, error)` | `SiteFromAddressIDForOrg`, `FeaturesFromSiteID` |
| **Select (global)** | `func XxxFromYyy(ctx, ...) (model, error)` | `CommunicationFromID`, `AddressFromID` |
| **Update (txn)** | `func XxxSetYyy(ctx, txn db.Ex, ...) error` | `CommunicationSetStatus` |
| **Bulk insert (txn)** | `func XxxInserts(ctx, txn db.Ex, []model) ([]model, error)` | `AddressInserts`, `ReportImagesInsert` |
| **Bulk select (txn)** | `func XxxsFromYyys(ctx, txn db.Ex, []int64) ([]model, error)` | `AddressesFromIDs`, `FeaturePoolsFromFeatures` |

After the signature normalization in CLEANUP.md §13, **every** function accepts `txn db.Ex`. All tests use the same transaction-based pattern: begin → insert → query → verify → rollback.

### The `db.Ex` interface (from `db/tx.go`)

```go
type Ex interface {
    Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
    Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
```

`db.BeginTxn()` returns `db.Tx` which implements `Ex`. `*pgxpool.Pool` does NOT implement `Ex` directly (it has different method signatures), which is why `db.ExecuteOne` (global pool) and `db.ExecuteOneTx` (transaction) are separate functions.

### Data flow

```
Query function constructs a Jet statement →
    calls db.ExecuteOneTx[T]() or db.ExecuteManyTx[T]() →
        statement.Sql() produces (query string, args) →
            txn.Query(ctx, query, args...) →
                pgx collects rows into typed struct
```

---

## Test Infrastructure

### Prerequisites

A running Postgres instance accessible via a connection string. The test framework will:

1. Connect using `TEST_POSTGRES_DSN` env var (default: skip tests if unset, so `go test ./...` works without DB)
2. Run all migrations via goose (embedded in `db/migrations/*.sql`)
3. Each test begins a transaction, runs the test, rolls back
4. No test data persists

### Test helper: `db/dbtest/dbtest.go`

Create a `db/dbtest/` package providing:

```go
package dbtest

import (
    "context"
    "os"
    "sync"
    "testing"

    "github.com/Gleipnir-Technology/nidus-sync/db"
    "github.com/jackc/pgx/v5/pgxpool"
)

var (
    pool *pgxpool.Pool
    once sync.Once
)

// Setup ensures the test database is initialized (migrations run).
// Called once per test binary via TestMain or per-package init.
func Setup(t *testing.T) {
    t.Helper()
    dsn := os.Getenv("TEST_POSTGRES_DSN")
    if dsn == "" {
        t.Skip("TEST_POSTGRES_DSN not set, skipping DB tests")
    }
    once.Do(func() {
        ctx := context.Background()
        if err := db.InitializeDatabase(ctx, dsn); err != nil {
            t.Fatalf("initialize test database: %v", err)
        }
        pool = db.PGInstance.PGXPool
    })
}

// Txn begins a new transaction on the test pool and returns
// it along with a rollback cleanup function.
func Txn(t *testing.T) (context.Context, db.Ex, func()) {
    t.Helper()
    ctx := context.Background()
    tx, err := pool.Begin(ctx)
    if err != nil {
        t.Fatalf("begin txn: %v", err)
    }
    return ctx, tx, func() {
        tx.Rollback(ctx)
    }
}

// Pool returns the raw pgxpool for tests that need it.
func Pool() *pgxpool.Pool {
    return pool
}
```

### Test file naming

All test files follow the standard Go convention: `<name>_test.go`, placed in the same package being tested (using `_test` external test package where needed for circular dependency avoidance). The package name follows `package public_test` pattern to test exported API only.

Actually, since the query functions are all exported and testable from outside, use:

```go
package public_test  // external test package
```

This avoids circular dependency on `db/dbtest` and ensures we only test the public API.

### Dependencies to add to `go.mod`

Promote to direct (already indirect):
```
github.com/stretchr/testify v1.11.1
```

Add for assertions:
```
require "github.com/stretchr/testify/require"
assert "github.com/stretchr/testify/assert"
```

---

## Phase 1: INSERT Functions (lowest risk, highest clarity)

These are the simplest: construct a model, insert, verify the returned model has an auto-generated ID.

### 1.1 `db/query/public/` insert functions

| File | Function | Model Dependencies | Notes |
|------|----------|-------------------|-------|
| `communication.go` | `CommunicationInsert` | `Communication` | Requires `organization_id` FK. Create an org in test setup. |
| `communication_log_entry.go` | `CommunicationLogEntryInsert` | `CommunicationLogEntry` | Requires `communication_id` FK. Insert a communication first. |
| `compliance_report_request.go` | `ComplianceReportRequestInsert` | `ComplianceReportRequest` | Requires `lead_id` FK (nullable). Test with nil. |
| `lead.go` | `LeadInsert` | `Lead` | Requires `organization_id` and `site_id` FK. |
| `signal.go` | `SignalInsert` | `Signal` | Requires `organization_id`, `location` (geom), FK to `site_id` (nullable). |
| `job.go` | `JobInsert` | `Job` | Verify FK constraints documented. |

### 1.2 `db/query/publicreport/` insert functions

| File | Function | Model Dependencies |
|------|----------|-------------------|
| `compliance.go` | `ComplianceInsert` | `Compliance` |
| `image.go` | `ImageInsert` | `Image` |
| `image_exif.go` | `ImageExifInserts` | `ImageExif` (bulk) |
| `nuisance.go` | `NuisanceInsert` | `Nuisance` |
| `report.go` | `ReportInsert` | `Report` |
| `report_image.go` | `ReportImageInsert`, `ReportImagesInsert` | `ReportImage` (single + bulk) |
| `report_log.go` | `ReportLogInsert` | `ReportLog` |
| `water.go` | `WaterInsert` | `Water` |

### 1.3 `db/query/arcgis/` insert functions

| File | Function | Model Dependencies |
|------|----------|-------------------|
| `account.go` | `AccountInsert` | `Account` |

### Example test: `db/query/public/communication_test.go`

```go
package public_test

import (
    "testing"
    "time"

    "github.com/Gleipnir-Technology/nidus-sync/db/dbtest"
    "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
    query "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCommunicationInsert(t *testing.T) {
    dbtest.Setup(t)
    ctx, txn, cleanup := dbtest.Txn(t)
    defer cleanup()

    comm := model.Communication{
        Created:        time.Now(),
        OrganizationID: 1,
        Status:         model.Communicationstatus_New,
    }
    result, err := query.CommunicationInsert(ctx, txn, comm)

    require.NoError(t, err)
    assert.NotZero(t, result.ID)
    assert.Equal(t, model.Communicationstatus_New, result.Status)
    assert.Equal(t, int32(1), result.OrganizationID)
}
```

### Test structure pattern

Every INSERT test follows this template:

1. **Arrange**: Create a model struct with required fields populated.
2. **Act**: Call the Insert function with a test transaction.
3. **Assert**: 
   - No error returned.
   - `result.ID` is non-zero (auto-generated primary key).
   - Inserted values match input values.
   - Timestamp fields are set (where applicable).

---

## Phase 2: SELECT Functions

These require data to already exist in the table. Each SELECT test inserts a row in the same transaction, then queries it back. After the signature normalization (CLEANUP.md §13), **all** SELECT functions accept `txn db.Ex` and use `ExecuteOneTx`/`ExecuteManyTx`.

### 2.1 `db/query/public/` select functions

| File | Function | Strategy |
|------|----------|----------|
| `address.go` | `AddressFromComplianceReportRequestID` | Insert address → query by report request ID |
| `address.go` | `AddressFromGID` | Insert address → query by GID |
| `address.go` | `AddressFromID` | Insert address → query by ID |
| `address.go` | `AddressesFromIDs` | Insert 2 addresses → query by IDs |
| `communication.go` | `CommunicationFromID` | Insert communication → query by ID |
| `communication.go` | `CommunicationsFromOrganization` | Insert 2 communications → query by org |
| `feature.go` | `FeaturesFromSiteID` | Insert site → feature → query |
| `feature.go` | `FeaturesFromSiteIDs` | Insert 2 sites + features → query |
| `feature_pool.go` | `FeaturePoolsFromFeatures` | Insert feature → pool → query |
| `site.go` | `SiteFromAddressIDForOrg` | Insert address + site → query |
| `site.go` | `SiteFromIDForOrg` | Insert site → query |

### 2.2 `db/query/publicreport/` select functions

| File | Function | Strategy |
|------|----------|----------|
| `compliance.go` | `ComplianceFromID` | Insert compliance → query by ID |
| `report.go` | `ReportFromPublicID` | Insert report → query by public ID |
| `report.go` | `ReportFromPublicIDForOrg` | Insert report → query by public ID + org |
| `report.go` | `ReportFromID` | Insert report → query by ID |
| `report.go` | `ReportsFromIDs` | Insert 2 reports → query by IDs |
| `report.go` | `ReportsFromIDsForOrg` | Insert 2 reports → query by IDs + org |
| `report.go` | `ReportsUnreviewedForOrganization` | Insert reviewed + unreviewed → query unreviewed |

### 2.3 `db/query/arcgis/` select functions

| File | Function | Strategy |
|------|----------|----------|
| `account.go` | `AccountFromID` | Insert account → query by ID |
| `oauth.go` | `OAuthTokenFromID` | Insert token → query by ID |
| `oauth.go` | `OAuthTokenForUser` | Insert token → query by user |
| `oauth.go` | `OAuthTokensForUser` | Insert tokens → query by user |
| `oauth.go` | `OAuthTokensValid` | Insert valid + invalid → query valid |
| `oauth.go` | `OAuthTokenForUserExists` | Insert token → verify exists |
| `service_feature.go` | `ServiceFeatureFromID` | Insert → query by ID |
| `service_feature.go` | `ServiceFeatureFromURL` | Insert → query by URL |
| `service_map.go` | `ServiceMapFromID` | Insert → query by ID |
| `service_map.go` | `ServiceMapsFromAccountID` | Insert maps → query by account |

### Example test: `db/query/public/address_test.go`

```go
package public_test

import (
    "testing"
    "time"

    "github.com/Gleipnir-Technology/nidus-sync/db/dbtest"
    "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
    query "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/twpayne/go-geom"
)

func TestAddressFromGID(t *testing.T) {
    dbtest.Setup(t)
    ctx, txn, cleanup := dbtest.Txn(t)
    defer cleanup()

    // Insert test data
    addr := model.Address{
        Country:    "US",
        Created:    time.Now(),
        Location:   geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{-122.4, 37.8}),
        H3cell:     "test",
        Locality:   "San Francisco",
        PostalCode: "94102",
        Street:     "Market St",
        Unit:       "",
        Region:     "CA",
        Number:     "1234",
        Gid:        "test-gid-001",
    }
    inserted, err := query.AddressInsert(ctx, txn, addr)
    require.NoError(t, err)

    // Query by GID
    result, err := query.AddressFromGID(ctx, txn, "test-gid-001")
    require.NoError(t, err)
    require.NotNil(t, result)
    assert.Equal(t, inserted.ID, result.ID)
    assert.Equal(t, "US", result.Country)
    assert.Equal(t, "San Francisco", result.Locality)
}

func TestAddressesFromIDs(t *testing.T) {
    dbtest.Setup(t)
    ctx, txn, cleanup := dbtest.Txn(t)
    defer cleanup()

    // Insert two addresses
    a1, _ := query.AddressInsert(ctx, txn, model.Address{
        Created: time.Now(), Location: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{0, 0}),
        H3cell: "a", Gid: "gid-a",
    })
    a2, _ := query.AddressInsert(ctx, txn, model.Address{
        Created: time.Now(), Location: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{0, 1}),
        H3cell: "b", Gid: "gid-b",
    })

    // Query by IDs
    results, err := query.AddressesFromIDs(ctx, txn, []int64{int64(a1.ID), int64(a2.ID)})
    require.NoError(t, err)
    assert.Len(t, results, 2)

    // Empty input returns empty output
    empty, err := query.AddressesFromIDs(ctx, txn, []int64{})
    require.NoError(t, err)
    assert.Empty(t, empty)
}
```

---

## Phase 3: UPDATE Functions

Verify that updates modify rows correctly and respect predicates.

### 3.1 Update functions

| File | Function | Test Strategy |
|------|----------|---------------|
| `communication.go` | `CommunicationSetStatus` | Insert communication → update status → verify via SELECT |
| `report.go` | `ReportUpdater.Execute` | Insert report → apply updater → verify |

### Example test: `db/query/public/communication_test.go`

```go
func TestCommunicationSetStatus(t *testing.T) {
    dbtest.Setup(t)
    ctx, txn, cleanup := dbtest.Txn(t)
    defer cleanup()

    // Insert
    comm, err := query.CommunicationInsert(ctx, txn, model.Communication{
        Created:        time.Now(),
        OrganizationID: 1,
        Status:         model.Communicationstatus_New,
    })
    require.NoError(t, err)

    // Update status
    err = query.CommunicationSetStatus(ctx, txn,
        int64(comm.OrganizationID), int64(comm.ID),
        model.Communicationstatus_Closed)
    require.NoError(t, err)

    // Verify the update via a SELECT in the same transaction
    // (CommunicationFromID accepts db.Ex after CLEANUP.md §13a conversion)
    updated, err := query.CommunicationFromID(ctx, txn, int64(comm.ID))
    require.NoError(t, err)
    assert.Equal(t, model.Communicationstatus_Closed, updated.Status)
}
```

---

## Phase 4: ArcGIS Query Package

After the `bob.Tx` → `db.Ex` conversion (CLEANUP.md §13c), the arcgis query functions use the same transaction pattern as the other packages.

### 4.1 INSERT functions

| File | Function | Notes |
|------|----------|-------|
| `account.go` | `AccountInsert` | After 13c: uses `ExecuteOneTx` |
| `oauth.go` | `OAuthTokenInsert` | After 13a: accepts `txn db.Ex` |
| `service_feature.go` | `ServiceFeatureInsert` | After 13c: uses `ExecuteOneTx` |
| `service_map.go` | `ServiceMapInsert` | After 13c: uses `ExecuteOneTx` |
| `user.go` | `UserInsert` | After 13c: uses `ExecuteOneTx` |
| `user_privileges.go` | `UserPrivilegeInsert` | After 13c: uses `ExecuteOneTx` |

### 4.2 SELECT functions

| File | Function | Notes |
|------|----------|-------|
| `account.go` | `AccountFromID` | After 13a: accepts `txn db.Ex` |
| `oauth.go` | `OAuthTokenFromID`, `OAuthTokenForUser`, `OAuthTokensForUser`, `OAuthTokensValid`, `OAuthTokenForUserExists` | After 13a |
| `service_feature.go` | `ServiceFeatureFromID`, `ServiceFeatureFromURL` | After 13a |
| `service_map.go` | `ServiceMapFromID`, `ServiceMapsFromAccountID` | After 13a |

### 4.3 UPDATE/DELETE functions

| File | Function | Notes |
|------|----------|-------|
| `oauth.go` | `OAuthTokenUpdateAccessToken`, `OAuthTokenUpdateRefreshToken`, `OAuthTokenUpdateLicense`, `OAuthTokenInvalidate` | After 13a |
| `user_privileges.go` | `UserPrivilegesDeleteByUserID` | After 13c |

---

## Phase 5: Edge Cases and Error Handling

### 5.1 Empty bulk operations

Functions like `AddressesFromIDs` and `ReportImagesInsert` already handle empty input slices gracefully. Write tests confirming:
- Empty input → non-nil empty slice, no error.
- Nil input → handled gracefully (or skipped with `t.Skip` if the function panics).

### 5.2 Unique constraint violations

Insert two rows with the same unique key; verify the error message is readable.

### 5.3 Foreign key violations

Insert a row referencing a non-existent parent; verify the error. This validates that FK constraints are correctly defined in the schema.

### 5.4 Not found

Functions returning `(*model.Xxx, error)` should return `nil, nil` on not-found (pattern already used by `ReportFromPublicID` and `AddressFromGID`). Test this behavior.

### 5.5 NULL handling

Models with nullable fields (`*int32`, `*string`, `*time.Time`, `*geom.T`): test with nil and non-nil values to verify round-trip fidelity.

---

## Test Execution

### Local development

```bash
# Set up a test database (one time)
createdb nidus-sync-test

# Run the query-layer tests
TEST_POSTGRES_DSN="postgresql://?host=/var/run/postgresql&dbname=nidus-sync-test" \
  go test ./db/query/... -v -count=1

# Run all tests (skips DB tests if no DSN set)
go test ./... -v -count=1
```

### CI (GitHub Actions example)

```yaml
services:
  postgres:
    image: postgres:16
    env:
      POSTGRES_DB: nidus-sync-test
      POSTGRES_PASSWORD: password
    options: >-
      --health-cmd pg_isready
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5
    ports:
      - 5432:5432

steps:
  - name: Test query layer
    run: |
      go test ./db/query/... -v -count=1
    env:
      TEST_POSTGRES_DSN: "postgresql://postgres:password@localhost:5432/nidus-sync-test"
```

### Using test databases in sequence vs parallel

- All Phase 1 INSERT tests can run in parallel (they use separate transactions on separate tables).
- All SELECT tests for the same table should run sequentially to avoid ID conflicts.
- Per-package `TestMain` can handle `db.InitializeDatabase` once, then run all tests.

**Recommended approach**: Run all tests sequentially within each package (Go's default), using `-count=1` to disable caching. Each test starts its own transaction, so there's no data leakage even running sequentially.

---

## File-by-File Implementation Order

### Step 1: Infrastructure (1 file)

| File | Purpose |
|------|---------|
| `db/dbtest/dbtest.go` | Test helper: pool setup, migration runner, transaction factory |

### Step 2: `go.mod` change (1 line)

Promote `github.com/stretchr/testify` to direct dependency.

### Step 3: INSERT tests (8 test files)

| Test File | Query File Tested | Functions Covered |
|-----------|------------------|-------------------|
| `db/query/public/communication_test.go` | `communication.go` + `communication_log_entry.go` | `CommunicationInsert`, `CommunicationLogEntryInsert`, `CommunicationSetStatus`, `CommunicationFromID`, `CommunicationsFromOrganization` |
| `db/query/public/address_test.go` | `address.go` | `AddressInsert`, `AddressesFromIDs`, `AddressFromGID`, `AddressFromID`, `AddressFromComplianceReportRequestID` |
| `db/query/public/site_test.go` | `site.go` | `SiteFromAddressIDForOrg`, `SiteFromIDForOrg` |
| `db/query/public/lead_test.go` | `lead.go` | `LeadInsert` |
| `db/query/public/signal_test.go` | `signal.go` | `SignalInsert` |
| `db/query/public/compliance_report_request_test.go` | `compliance_report_request.go` | `ComplianceReportRequestInsert` |
| `db/query/public/feature_test.go` | `feature.go` + `feature_pool.go` | `FeaturesFromSiteID`, `FeaturePoolsFromFeatures`, `FeaturesFromSiteIDs` |
| `db/query/publicreport/report_test.go` | `report.go` | `ReportInsert`, `ReportFromPublicID`, `ReportFromPublicIDForOrg`, `ReportFromID`, `ReportsFromIDs`, `ReportsFromIDsForOrg`, `ReportsUnreviewedForOrganization` |

### Step 4: Remaining query packages (4 test files)

| Test File | Query File Tested | Functions Covered |
|-----------|------------------|-------------------|
| `db/query/publicreport/compliance_test.go` | `compliance.go` | `ComplianceInsert`, `ComplianceFromID` |
| `db/query/publicreport/image_test.go` | `image.go` + `image_exif.go` + `report_image.go` | All image insert functions |
| `db/query/publicreport/nuisance_water_test.go` | `nuisance.go` + `water.go` + `report_log.go` | `NuisanceInsert`, `WaterInsert`, `ReportLogInsert` |
| `db/query/arcgis/arcgis_test.go` | `account.go` + `oauth.go` + `service_feature.go` + `service_map.go` + `user.go` | All arcgis query functions (after 13a + 13c conversions) |

---

## Model Foreign Key Dependency Graph

Understanding which inserts require which parent rows helps plan test setup:

```
organization ─────────────────────────────────────────────┐
    │                                                      │
    ├── communication ── communication_log_entry            │
    ├── site ── feature ── feature_pool                     │
    │     │                                                 │
    │     ├── signal (site_id, location)                    │
    │     └── lead (site_id) ── compliance_report_request   │
    │                                                       │
    └── publicreport.report ── report_log                   │
           ├── report_image                                │
           ├── compliance (report_id)                      │
           ├── nuisance (report_id)                        │
           └── water (report_id)                           │
```

For initial INSERT tests, we need at minimum a test `organization` row. The `dbtest.Setup` function can optionally seed this.

### Seeding approach

Option A — Seed in `dbtest.Setup()`: insert a minimal org row (id=1) during migration/setup so all tests have a valid FK target.
Option B — Each test creates its own dependency rows within the transaction.

**Recommendation**: Option B for now (each test is self-contained). The overhead is low and tests remain independent. If organization-schema evolves and gets more columns, we can add a helper:

```go
func SeedOrganization(ctx context.Context, txn db.Ex) (int32, error) {
    // Insert a minimal org row
}
```

---

## What Is NOT Covered (yet)

| Area | Reason |
|------|--------|
| `db/prepared.go` param builders | Scheduled for removal (per project owner) |
| Platform layer (`platform/*.go`) | Separate plan — these call query functions; test them after query layer is solid |
| HTTP handlers (`api/`, `resource/`) | Need HTTP test infrastructure (httptest) |
| Bob ORM-generated models (`db/models/`) | Legacy ORM; query tests cover the Jet layer which is the migration target |
| `db/fieldseeker.go` | Entirely commented out |
| `db/connection.go` `Execute*` helpers | Covered transitively by query tests; direct tests would be lower priority |
| Vue/TypeScript frontend | Separate test effort (Vitest) |

---

## Success Criteria

After all phases complete:

1. **Signature normalization (CLEANUP.md §13)**: every query function has `(ctx context.Context, txn db.Ex, ...)` signature. No function uses the global pool internally.
2. **Every exported function in `db/query/public/`**, `db/query/publicreport/`, and `db/query/arcgis/` has at least one transaction-based test.
3. **INSERT functions**: verify returned model has auto-generated ID and correct typed fields.
4. **SELECT functions**: verify round-trip (insert → query → match) within the same transaction.
5. **UPDATE functions**: verify modification takes effect, verifiable via SELECT in same transaction.
6. **Edge cases**: empty input slices, not-found returns `nil`/error, FK/unique violations produce errors, NULL round-trips.
7. **CI green**: tests pass in CI with a Postgres service container.
