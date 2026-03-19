# Proposal: Normalize Numeric `object_id` Fields

## Why

The provider currently exposes `object_id` as a string on generated resources
such as `awx_role_team_assignment`, even when the field semantically carries an
AWX numeric primary key. That breaks the provider's own numeric ID contract and
forces unnecessary type conversions when wiring role assignments to other AWX
objects.

## What Changes

- Define a provider rule that generated `object_id` fields representing AWX
  numeric primary keys use Terraform `Number` rather than `String`.
- Apply that rule to both configurable resource arguments and computed/read-only
  attributes for affected resources and data sources.
- Keep UUID-style alternatives such as `object_ansible_id` typed as strings.
- Preserve import identifiers and resource identity semantics; only field typing
  changes.
- Mark the schema type correction as **BREAKING** for configurations or state
  that currently treat affected `object_id` values as quoted strings.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `awx-typed-reference-ids`: Extend numeric reference typing requirements to
  cover semantic numeric `object_id` fields, not only canonical `*_id`
  references.
- `awx-single-object-resource-model`: Require generated managed resources to
  expose semantic numeric `object_id` fields as Terraform numbers in both plan
  input and state.
- `awx-data-sources-and-lookups`: Require generated data sources to expose
  semantic numeric `object_id` attributes as Terraform numbers when the
  underlying AWX identifier is numeric.

## Impact

- Generator and manifest typing logic in
  `/Users/damien/git/terraform-provider-awx-iwd/internal/openapi` and
  `/Users/damien/git/terraform-provider-awx-iwd/internal/manifest`.
- Generated provider schemas and docs for resources and data sources that
  expose numeric `object_id` fields, including
  `/Users/damien/git/terraform-provider-awx-iwd/docs/resources/awx_role_team_assignment.md`
  and related generated outputs.
- Provider tests covering generated manifest field types, resource schemas, and
  data source schemas.
- User configurations and state that currently quote affected `object_id`
  values.
