# Tasks: Fix Organization Credential Relationship Mapping

## 1. Relationship Derivation Fix

- [x] 1.1 Add child collection alias resolution for `galaxy_credentials` ->
  `credentials` in `internal/openapi/parser.go`.
- [x] 1.2 Implement deterministic collision preference so
  `organization_credential_association` resolves to
  `/api/v2/organizations/{id}/galaxy_credentials/` when both credential paths
  exist.
- [x] 1.3 Ensure relationship naming and canonical argument naming remain
  unchanged (`organization_id`, `credential_id`).

## 2. Regression Tests

- [x] 2.1 Add parser unit test covering dual organization credential paths and
  asserting selected path is `.../galaxy_credentials/`.
- [x] 2.2 Add parser non-regression assertions that existing special variants
  (`notification_templates_*`, `survey_spec`) are unaffected.
- [x] 2.3 Keep/verify client test asserting association requests use
  `.../galaxy_credentials/` for organization credential relationships.

## 3. Generation and Validation

- [x] 3.1 Run `make generate`.
- [x] 3.2 Run `make validate-manifest`.
- [x] 3.3 Run `make docs`.
- [x] 3.4 Run `make docs-validate`.
- [x] 3.5 Run `make test`.

## 4. Acceptance and Evidence

- [x] 4.1 Verify generated
  `internal/manifest/relationships.json` maps
  `organization_credential_association` to
  `/api/v2/organizations/{id}/galaxy_credentials/`.
- [x] 4.2 Document test evidence and notable manifest/docs diffs in the change
  implementation summary (or commit/PR notes).
