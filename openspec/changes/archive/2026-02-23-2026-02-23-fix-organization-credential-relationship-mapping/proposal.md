# Proposal: Fix Organization Credential Relationship Mapping

## Why

`awx_organization_credential_association` currently maps to
`/api/v2/organizations/{id}/credentials/`, which AWX treats as credential
creation semantics in practice for this flow, leading to 400 errors requiring
creation fields (`name`, `credential_type`) instead of accepting attach-by-id.

This indicates a relationship-derivation ambiguity in provider generation:
multiple organization credential-related paths exist (`credentials` and
`galaxy_credentials`), but the current path selection picks the wrong endpoint
for association lifecycle behavior.

## What Changes

- Adjust relationship path derivation to map organization credential association
  to `/api/v2/organizations/{id}/galaxy_credentials/`.
- Add deterministic collision handling when multiple paths map to the same
  relationship name so endpoint selection is intentional, not lexical-order
  accidental.
- Add parser-level and client-level non-regression coverage.
- Regenerate manifests/docs and validate consistency.

## Capabilities

### Modified Capabilities

- `awx-relationship-resources`: relationship derivation must select endpoint
  paths that preserve association semantics (attach/detach) when aliases or
  overlapping child collections exist.

## Impact

- Generator logic in `internal/openapi/parser.go`.
- Relationship metadata in `internal/manifest/relationships.json` (generated).
- Client safeguard logic in `internal/client/client.go` and tests.
- Relationship behavior documentation generated from manifests.
