## Why

Reference fields that point to other AWX objects are currently named like plain attributes (for example, `organization`), which makes configuration intent ambiguous. Adding an explicit `_id` suffix for link fields improves readability and makes resource and data source wiring behavior clear at a glance.

## What Changes

- Rename Terraform schema fields that represent AWX object links to use the `_id` suffix (for example, `organization` -> `organization_id`).
- Apply the naming contract consistently for configurable arguments and computed attributes where the value is an object identifier reference.
- **BREAKING**: Remove or deprecate unsuffixed reference field names from generated resources/data sources according to the migration strategy defined in design/specs.
- Regenerate provider documentation/examples so link wiring always uses explicit `_id` field names.
- Add tests that verify link-field naming consistency across schema generation, provider runtime behavior, and docs output.

## Capabilities

### New Capabilities
- `awx-reference-id-field-naming`: Defines the provider-wide naming contract that link/reference fields must use `_id` suffixes so users can clearly distinguish object links from non-reference fields.

### Modified Capabilities
- `awx-single-object-resource-model`: Updates object schema requirements so reference fields are surfaced with explicit `_id` names.
- `awx-data-sources-and-lookups`: Aligns exported lookup/reference attributes with the `_id` naming contract.
- `awx-provider-documentation-and-examples`: Updates argument/attribute docs and examples to demonstrate link wiring through `_id` fields.
- `awx-typed-reference-ids`: Clarifies that typed reference-ID behavior applies to suffixed reference fields and not ambiguous unsuffixed names.

## Impact

- Affected code: schema/manifest derivation and field mapping in `internal/openapi/*`, `internal/manifest/*`, and resource/data source schema assembly in `internal/provider/*`.
- Affected outputs: generated manifests and docs, including `internal/manifest/*.json` and `docs/resources/*.md` / `docs/data-sources/*.md`.
- User impact: clearer configuration semantics for resource linking, with migration required for configurations using unsuffixed reference field names.
- Validation impact: unit/provider/acceptance coverage must assert naming consistency and migration behavior.
