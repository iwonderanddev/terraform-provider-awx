# AGENTS.md

## Purpose

This repository implements a Terraform provider for AWX API v2, with broad coverage driven by a vendored OpenAPI schema plus curated metadata overrides.  
Target compatibility is AWX `24.6.1` (API `/api/v2`), with HTTP Basic authentication at GA.

## Quick Facts

- Language: Go `1.25.0`
- Provider address: `registry.terraform.io/iwd/awx`
- Resource/data source registration is dynamic from embedded manifest JSON
- Object import IDs:
  - Collection-created objects: numeric (`42`)
  - Detail-path keyed objects (for example settings): detail identifier (`system`)
- Relationship import IDs:
  - Standard associations: `<primary_id>:<related_id>` (for example `12:34`)
  - Survey spec relationships: `<resource_id>` (for example `12`)

## Repository Map

- `cmd/terraform-provider-awx-iwd/main.go`: provider server entrypoint
- `cmd/awxgen/main.go`: generator/validator/docs CLI
- `internal/provider/*`: provider runtime for objects, relationships, data sources
- `internal/client/*`: AWX HTTP transport, retries, pagination, errors
- `internal/openapi/*`: schema parsing + manifest derivation logic
- `internal/manifest/*`: curated controls + generated manifests (embedded at runtime)
- `internal/acceptance/*`: direct client-level live AWX acceptance tests
- `docs/*`: generated provider/resource/data-source docs
- `examples/*`: Terraform usage examples
- `external/awx-openapi/schema.json`: vendored AWX OpenAPI source schema

## Source Of Truth Rules

### Curated (manual edits expected)

- `internal/manifest/runtime_exclusions.json`
- `internal/manifest/deprecated_exclusions.json`
- `internal/manifest/relationship_priorities.json`
- `internal/manifest/field_overrides.json`
- `external/awx-openapi/schema.json` (normally via update script)
- `external/awx-openapi/README.md` (retrieval date/hash after schema updates)

### Generated (do not hand-edit)

- `internal/manifest/managed_objects.json`
- `internal/manifest/relationships.json`
- `internal/manifest/coverage_report.json`
- `docs/index.md`
- `docs/resources/*.md`
- `docs/data-sources/*.md`

After changing curated inputs or schema, always run:

1. `make generate`
2. `make validate-manifest`
3. `make docs`
4. `make docs-validate`
5. `make test`
6. `make build`

## Core Architecture

1. `awxgen` loads `external/awx-openapi/schema.json`.
2. `internal/openapi` derives managed objects and relationship candidates.
3. Curated exclusions/priorities/field overrides/deprecation exclusions are applied.
4. Generated manifests are written to `internal/manifest/*.json`.
5. `internal/manifest.Load()` embeds and loads catalog metadata at provider startup.
6. Provider dynamically registers:
   - object resources (`object_resource.go`)
   - relationship resources (`relationship_resource.go`)
   - object data sources (`object_data_source.go`)

## Runtime Behavior Details

### Object resources

- `collectionCreate=true`:
  - Create via `POST` collection endpoint
  - State `id` is numeric
- `collectionCreate=false`:
  - Lifecycle keyed by detail-path identifier from config (`id` required)
  - Non-numeric IDs are allowed for these resources
- Write-only sensitive fields are sent on create/update and preserved from plan/state; they are not read back from API responses.
- Array/object fields are represented in Terraform as JSON-encoded strings.
- Optional fields with AWX OpenAPI defaults are often marked `Optional + Computed` to avoid null/default drift after apply.

### Relationship resources

- Standard relationships:
  - Inputs: canonical object-specific `*_id` attributes (for example `team_id`, `user_id`)
  - State/import ID: `<primary_id>:<related_id>`
- Survey spec relationships (`.../survey_spec/`) are special:
  - Inputs: canonical parent object `*_id` attribute + `spec` (JSON string)
  - State/import ID: `<resource_id>`

### Data sources

- Deterministic lookup order:
  1. `id` if provided
  2. `name` exact match (only if the object has a `name` field)
- Name lookup must resolve exactly one result; 0 or >1 is an error.

## AWX Client Contract

- Auth: HTTP Basic only (`username` + `password`)
- TLS options: `insecure_skip_tls_verify`, `ca_cert_pem`
- Retry behavior with backoff for retryable failures
- Pagination support via `ListAll` following AWX `next` links
- `Ping()` validates connectivity against `/api/v2/` during provider configure

## Commands

- `make generate`: regenerate manifests + coverage report
- `make validate-manifest`: verify manifests are in sync with generator output
- `make docs`: regenerate docs from manifests
- `make docs-validate`: verify docs presence and required sections
- `make coverage-report`: print coverage summary
- `make test`: run all Go tests
- `make test-acceptance`: run opt-in live AWX tests (`internal/acceptance` + Terraform-driven provider acceptance tests)
- `make build`: build `dist/terraform-provider-awx-iwd`

## Markdown Quality Gate

- When creating or editing any `*.md` file, use `$markdownlint-auto-fix`; the skill is responsible for running markdownlint, applying `--fix`, resolving remaining violations, and finishing only when lint returns zero errors.

## Acceptance Testing

`.env` is loaded by `make test-acceptance` if present.

Required:

- `AWX_ACCEPTANCE=1`
- `AWX_BASE_URL`
- `AWX_USERNAME`
- `AWX_PASSWORD`

Scenario fixture vars:

- `AWX_TEST_ORGANIZATION_ID`
- `AWX_TEST_TEAM_ID`
- `AWX_TEST_USER_ID`

Suites:

- `internal/acceptance`: direct client-level live API behavior
- `internal/provider`: Terraform lifecycle/import checks via `terraform-plugin-testing`

## Common Change Playbooks

### 1) Refresh AWX schema

1. `./external/awx-openapi/update.sh`
2. Update date/hash in `external/awx-openapi/README.md`
3. Run generate/validate/docs/test command chain

### 2) Runtime-only object appears in coverage validation

1. Confirm object is runtime-only/non-desired-state
2. Add exclusion to `internal/manifest/runtime_exclusions.json`
3. Regenerate + validate + docs + tests

### 3) Schema/runtime mismatch for a field

1. Add/adjust entry in `internal/manifest/field_overrides.json`
2. Regenerate + validate manifests
3. Update docs and tests as needed

### 4) Relationship ordering/name behavior needs adjustment

1. Modify `internal/manifest/relationship_priorities.json`
2. Regenerate and validate manifests/docs/tests

### 5) Deprecated endpoint should be removed from provider surface

1. Add object/path entry to `internal/manifest/deprecated_exclusions.json`
2. Regenerate and validate manifests/docs/tests

## Invariants To Preserve

- Keep AWX-native naming alignment (`awx_<singular>` and generated relationship resource names).
- Preserve import ID contracts (numeric/detail-key for objects; composite/parent-key for relationships).
- Do not repopulate write-only secret values from read responses.
- Keep generated manifests and docs in sync before concluding a change.
- Treat runtime-only objects via explicit exclusions, not ad hoc runtime branching.

## Current Coverage Snapshot

From `internal/manifest/coverage_report.json` (generated `2026-02-10`):

- `39` total object candidates
- `23` managed object resources
- `24` managed object data sources
- `14` runtime exclusions
- `61` relationship resources
