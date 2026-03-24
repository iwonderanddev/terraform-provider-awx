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
make build

```

Configure Terraform CLI dev override in `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/iwd/awx" = "/Users/damien/git/terraform-provider-awx-iwd/dist"
  }
  direct {}
}
```

Create a minimal Terraform config (example):

```hcl
terraform {
  required_providers {
    awx = {
      source = "iwd/awx"
    }
  }
}

provider "awx" {
  hostname = "https://awx.example.com"
  username = "admin"
  password = "..."
}
```

Then run:

```bash
# with dev_overrides, build the provider binary and skip terraform init
make build
terraform plan
```

## Compatibility

- AWX: `24.6.1`
- API: `/api/v2`
- Auth at GA: HTTP Basic (`username` + `password`)
- Object import IDs: numeric (`42`)
- Relationship import IDs: composite (`<primary_id>:<related_id>`, for example `12:34`)

## Generated Docs Qualifiers

Generated resource docs in `docs/resources/*` use these argument qualifiers:

- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX may apply a server-side default and Terraform stores the resulting value in state after apply.

`Optional, Computed` is used to avoid apply inconsistencies when AWX returns defaults for unset fields.

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
- `awx_team_user_association` relationship lifecycle + composite import id (`<primary_id>:<related_id>`).

Run only Terraform-driven acceptance tests:

```bash
set -a; . ./.env; set +a
TF_ACC=1 go test ./internal/provider -run TestAcceptanceTerraform -v
```

If tests show `SKIP`, check `.env` values and confirm `AWX_ACCEPTANCE=1`.

## Build Provider Binary

```bash
make build
```

`make build` injects a provider version at build time:

- If `HEAD` is tagged (example `v0.2.3`), that exact tag is used.
- If `HEAD` is not tagged, a dev version is used (example `v0.2.3-dev.a1b2c3d`).
- If the workspace has uncommitted changes, `.dirty` is appended.

Override the injected version manually:

```bash
make build VERSION=v0.2.3
```

## Create A New Version (Developers)

Use SemVer tags (`vMAJOR.MINOR.PATCH`) to create release versions.

1. Ensure the branch is ready and tests pass:

```bash
make test
```

1. Create an annotated tag for the new version:

```bash
git tag -a v0.2.3 -m "Release v0.2.3"
```

1. Push the tag:

```bash
git push origin v0.2.3
```

1. Build from that tagged commit:

```bash
make build
```

The resulting binary reports version `v0.2.3` to Terraform.

## Publish to Terraform Registry

Terraform Registry lists versions only from **signed [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases)** on the public repository linked when you [publish the provider](https://developer.hashicorp.com/terraform/registry/providers/publishing). A bare Git tag or GitLab-only tag is not enough.

### One-time setup

1. Add your **GPG public key** (ASCII-armored) under [Terraform Registry → User Settings → Signing Keys](https://registry.terraform.io/settings/gpg-keys). Use an RSA or DSA key; the registry does not accept the default ECC type for signing.
1. In the **GitHub** repository: Settings → Secrets and variables → Actions, add:
   - `GPG_PRIVATE_KEY` — ASCII-armored secret key matching the registered public key
   - `PASSPHRASE` — passphrase for that key (if empty, use an empty secret or adjust the workflow)
1. Ensure [GitHub Actions](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository) are allowed to run (organization and repo settings).

### Cutting a release

1. Merge your release commit to the default branch and push it to **GitHub** (directly or via the GitLab `mirror_to_github` job).
1. Create and push a SemVer tag with a `v` prefix (for example `v1.0.0`) to **GitHub**. If GitLab is the source of truth, push the tag to GitLab; the `mirror_tag_to_github` CI job pushes that tag to GitHub so the tag exists where GitHub Actions runs.
1. Confirm the **Release** workflow (`.github/workflows/release.yml`) completes and publishes a non-draft release with ZIPs, `terraform-provider-awx_<version>_SHA256SUMS`, `terraform-provider-awx_<version>_SHA256SUMS.sig`, and the manifest JSON asset.
1. On the [provider page](https://registry.terraform.io/), use **Resync** if a new version does not appear (after fixing webhooks or the first successful release).

Local dry-run (optional): install [GoReleaser](https://goreleaser.com/install/) and run `goreleaser release --snapshot --clean` to verify builds without publishing.

## Use Locally With Terraform (Dev Override)

Because this provider is local development, use a Terraform CLI dev override.

`~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/iwd/awx" = "/Users/damien/git/terraform-provider-awx-iwd/dist"
  }
  direct {}
}
```

Then in a Terraform config:

```hcl
terraform {
  required_providers {
    awx = {
      source = "iwd/awx"
    }
  }
}

provider "awx" {
  hostname = "https://awx.example.com"
  username = "admin"
  password = "..."
}
```

Run:

```bash
make build
terraform plan
terraform apply
```

Important with `dev_overrides`: Terraform may still try to query the public
registry during `terraform init` and fail for unreleased namespaces (for
example `iwd/awx`). In local provider development, build the binary and run
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
