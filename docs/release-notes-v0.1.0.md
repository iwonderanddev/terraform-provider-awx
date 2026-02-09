# AWX Provider v0.1.0 (Draft)

## Support Contract

- Target: AWX 24.6.1
- API surface: `/api/v2`
- Authentication: HTTP Basic (`username` + `password`)
- Terraform import IDs:
  - Object resources: numeric AWX IDs (`42`)
  - Relationship resources: composite IDs (`<parent_id>:<child_id>`)

## Coverage and Exclusions

- Managed object candidates are generated from `external/awx-openapi/schema.json`.
- Runtime-only exclusions are explicitly tracked in `internal/manifest/runtime_exclusions.json`.
- Relationship resources are generated from supported parent/child association endpoints and prioritized by `internal/manifest/relationship_priorities.json`.

## Sensitive Handling

- Sensitive/password-like fields are marked sensitive in generated schemas.
- Write-only sensitive fields are preserved from configuration and not repopulated from read responses.

## Test Execution Guidance

- Unit tests: `GOCACHE=/tmp/go-build go test ./...`
- Manifest validation: `GOCACHE=/tmp/go-build go run ./cmd/awxgen validate`
- Acceptance/e2e (opt-in):
  - Set `AWX_ACCEPTANCE=1`
  - Set `AWX_BASE_URL`, `AWX_USERNAME`, `AWX_PASSWORD`
  - Team CRUD scenario requires `AWX_TEST_ORGANIZATION_ID`
  - Relationship scenario requires `AWX_TEST_TEAM_ID`, `AWX_TEST_USER_ID`

## Validation Results (2026-02-09)

- AWX target verified: `24.6.1` (`/api/v2/config/`)
- Manifest/doc validation:
  - `go run ./cmd/awxgen validate` -> pass
  - `go run ./cmd/awxgen docs-validate` -> pass
  - `go run ./cmd/awxgen report` -> 38 candidates, 20 resource-eligible, 43 relationship resources, 0 missing exclusions
- Unit/integration suite:
  - `go test ./...` -> pass
- Live acceptance/e2e against AWX 24.6.1:
  - `go test ./internal/acceptance -run TestAcceptance -v -count=1` -> pass
  - Executed scenarios:
    - Team object CRUD, import-id validation, and remote-delete verification
    - Team-user relationship create/read/delete with composite import-id semantics (`<parent_id>:<child_id>`)

## Known Limitations

- Runtime-only object exclusions are rule-based and may require curation updates as AWX evolves.
- Acceptance/e2e tests require an operator-provided AWX 24.6.1 environment and are intentionally skipped by default.
- Current live acceptance scenarios are representative (team + team-user association) and do not exhaustively CRUD every generated object category in one run.
