# Proposal: Rename Project Credential ID To SCM Credential ID

## Why

`awx_project` currently exposes the AWX project source-control credential as
`credential_id`. Although technically correct, that name is ambiguous in
Terraform because projects also interact with other credential concepts and the
resource already uses more precise names such as
`signature_validation_credential_id`.

This should be clarified at the schema level, not only in documentation, so
Terraform configuration directly reflects that the field is specifically the
SCM credential used to access source control content.

## What Changes

- **BREAKING** Rename `awx_project.credential_id` to
  `awx_project.scm_credential_id`.
- **BREAKING** Rename the corresponding computed data source attribute from
  `credential_id` to `scm_credential_id`.
- Remove the legacy `credential_id` field from the generated `awx_project`
  Terraform schema instead of keeping a compatibility alias.
- Add a declarative generator mechanism for object-field-specific canonical
  Terraform reference names so this rename is encoded in metadata rather than
  hard-coded in provider runtime logic.
- Regenerate manifests and provider docs/examples so generated output uses the
  SCM-specific field name consistently.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `awx-reference-id-field-naming`: generated object-link reference fields may
  require canonical semantic prefixes beyond the raw AWX field name when the
  generic name is ambiguous, and `awx_project` must expose
  `scm_credential_id`.
- `awx-provider-documentation-and-examples`: generated docs and examples for
  `awx_project` must describe and demonstrate the SCM-specific field name
  consistently for both resources and data sources.

## Impact

- Generator naming metadata and field-name derivation used for managed objects.
- Generated manifest output for `awx_project`.
- Provider resource/data source schema generated from embedded manifests.
- Generated docs and examples for `awx_project`.
- Acceptance and unit tests that assert project field names or import/state
  expectations.
