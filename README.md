# Terraform Provider for AWX

Terraform provider for AWX API v2, currently targeted at **AWX 24.6.1**.

This repo is not just handwritten resources. It uses AWX OpenAPI + generation metadata to produce and validate provider coverage.

## Quickstart (First 5 Minutes)

From the repo root:

```bash
cp .env.example .env
# edit .env with your AWX values

make generate
make validate-manifest
make test
make test-acceptance

go build -o terraform-provider-awx ./cmd/terraform-provider-awx
```

Configure Terraform CLI dev override in `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/damien/awx" = "/Users/damien/git/terraform-awx-provider"
  }
  direct {}
}
```

Create a minimal Terraform config (example):

```hcl
terraform {
  required_providers {
    awx = {
      source = "damien/awx"
    }
  }
}

provider "awx" {
  base_url = "https://awx.example.com"
  username = "admin"
  password = "..."
}
```

Then run:

```bash
# with dev_overrides, build the provider binary and skip terraform init
go build -o terraform-provider-awx ./cmd/terraform-provider-awx
terraform plan
```

## Compatibility

- AWX: `24.6.1`
- API: `/api/v2`
- Auth at GA: HTTP Basic (`username` + `password`)
- Object import IDs: numeric (`42`)
- Relationship import IDs: composite (`<parent_id>:<child_id>`, for example `12:34`)

## How This Project Works

High-level flow:

1. AWX OpenAPI schema is vendored at `external/awx-openapi/schema.json`.
2. Generator (`cmd/awxgen`) derives object/relationship metadata.
3. Generated metadata ("manifests") is stored under `internal/manifest/*.json`.
4. Provider runtime loads those manifests and dynamically registers resources/data sources.
5. Validation commands ensure generated output is in sync and docs are complete.

## What Are "Manifests"?

In this repo, **manifest** means generated JSON metadata files, not Terraform `.tf` files.

Main manifest files:

- `internal/manifest/managed_objects.json`
  - Derived object catalog from OpenAPI.
  - Includes resource/data-source eligibility and field metadata.
- `internal/manifest/relationships.json`
  - Derived association resource catalog.
- `internal/manifest/runtime_exclusions.json`
  - Explicit runtime-only object exclusions (curated input).
- `internal/manifest/relationship_priorities.json`
  - Priority ordering for relationship resource generation (curated input).
- `internal/manifest/field_overrides.json`
  - Manual field override hooks for schema/runtime mismatches (curated input).
- `internal/manifest/coverage_report.json`
  - Generated coverage summary used by validation.

If you change schema inputs or curated files, run `make generate`, then `make validate-manifest`.

## Repository Layout

- `cmd/terraform-provider-awx` provider server entrypoint
- `internal/provider` resource/data source/provider runtime
- `internal/client` shared AWX HTTP transport and error handling
- `internal/openapi` schema parsing and derivation logic
- `internal/manifest` generated metadata + curated controls
- `cmd/awxgen` generator and validation CLI
- `docs/` generated provider/resource/data-source docs
- `examples/` usage examples
- `internal/acceptance` opt-in live AWX acceptance/e2e tests

## Prerequisites

- Go (matching `go.mod`)
- AWX 24.6.1 reachable environment (for acceptance/e2e)

## Development Commands

- `make generate` regenerate manifest files from OpenAPI and curated controls
- `make validate-manifest` fail if committed manifests diverge from generated output
- `make docs` regenerate docs from manifests
- `make docs-validate` verify docs structure/completeness
- `make coverage-report` print coverage summary
- `make test` run all Go tests
- `make test-acceptance` run opt-in live AWX acceptance tests (loads `.env` if present), including Terraform-driven tests via `terraform-plugin-testing`

Typical local loop:

1. Edit provider/client/generator code or curated manifest inputs.
2. Run `make generate`.
3. Run `make validate-manifest`.
4. Run `make docs` (if manifest/schema changes).
5. Run `make test`.

## Acceptance / E2E (Opt-In)

Create local env file:

```bash
cp .env.example .env
```

Minimum `.env` values:

- `AWX_ACCEPTANCE=1`
- `AWX_BASE_URL`
- `AWX_USERNAME`
- `AWX_PASSWORD`

Additional scenario fixtures:

- `AWX_TEST_ORGANIZATION_ID`
- `AWX_TEST_TEAM_ID`
- `AWX_TEST_USER_ID`

Run:

```bash
make test-acceptance
```

### Terraform-Driven Acceptance Tests

`make test-acceptance` runs two acceptance suites:

- `internal/acceptance`: client-level live API checks (no Terraform runtime).
- `internal/provider`: Terraform-driven checks via `terraform-plugin-testing` (real `plan/apply/import` flows against AWX).

Terraform-driven scenarios currently cover:

- `awx_team` resource CRUD + import.
- `awx_team` data source lookup by name and id consistency checks.
- `awx_team_user_association` relationship lifecycle + composite import id (`<parent_id>:<child_id>`).

Run only Terraform-driven acceptance tests:

```bash
set -a; . ./.env; set +a
TF_ACC=1 go test ./internal/provider -run TestAcceptanceTerraform -v
```

If tests show `SKIP`, check `.env` values and confirm `AWX_ACCEPTANCE=1`.

## Build Provider Binary

```bash
go build -o terraform-provider-awx ./cmd/terraform-provider-awx
```

## Use Locally With Terraform (Dev Override)

Because this provider is local development, use a Terraform CLI dev override.

`~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/damien/awx" = "/Users/damien/git/terraform-awx-provider"
  }
  direct {}
}
```

Then in a Terraform config:

```hcl
terraform {
  required_providers {
    awx = {
      source = "damien/awx"
    }
  }
}

provider "awx" {
  base_url = "https://awx.example.com"
  username = "admin"
  password = "..."
}
```

Run:

```bash
go build -o terraform-provider-awx ./cmd/terraform-provider-awx
terraform plan
terraform apply
```

Important with `dev_overrides`: Terraform may still try to query the public
registry during `terraform init` and fail for unreleased namespaces (for
example `damien/awx`). In local provider development, build the binary and run
`terraform plan`/`terraform apply` directly from your test directory instead of
running `terraform init`.

## Updating Vendored AWX OpenAPI

Use:

```bash
./external/awx-openapi/update.sh
```

Then regenerate and validate:

```bash
make generate
make validate-manifest
make docs
make test
```
