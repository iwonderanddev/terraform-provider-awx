# Terraform Provider for AWX

This repository contains a Terraform provider targeting AWX 24.6.1 API v2.

## Development

- Generate manifests: `make generate`
- Validate manifest and generated mappings: `make validate-manifest`
- Generate docs: `make docs`
- Validate docs completeness: `make docs-validate`
- Run unit tests: `make test`

## Acceptance / E2E (Opt-In)

Acceptance tests are local and opt-in.

Required environment variables:

- `AWX_ACCEPTANCE=1`
- `AWX_BASE_URL`
- `AWX_USERNAME`
- `AWX_PASSWORD`

Additional scenario variables:

- `AWX_TEST_ORGANIZATION_ID` for object CRUD/import acceptance coverage
- `AWX_TEST_TEAM_ID` and `AWX_TEST_USER_ID` for relationship acceptance coverage

Run:

```bash
make test-acceptance
```

## Compatibility

- AWX 24.6.1 API v2 only
- HTTP Basic authentication at GA
- Runtime-only objects excluded via explicit manifest
