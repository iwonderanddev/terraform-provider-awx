## Why

Relationship resources currently use generic `parent_id` and `child_id` arguments, which makes Terraform configuration ambiguous and forces users to infer endpoint directionality. Introducing explicit relationship-specific ID arguments (for example, `job_template_id` and `credential_id`) makes link configuration self-documenting and reduces avoidable user errors.

## What Changes

- Define canonical relationship resource arguments using explicit object names with `_id` suffixes instead of generic `parent_id`/`child_id`.
- Update relationship resource generation/runtime mapping so each association exposes deterministic, endpoint-specific ID argument names.
- Preserve existing relationship import identity formats (`<parent_id>:<child_id>` and singleton `<parent_id>`) to avoid state/import contract churn.
- **BREAKING**: Remove `parent_id`/`child_id` from relationship resource schemas and require canonical object-specific `*_id` arguments.
- Regenerate relationship documentation/examples to demonstrate explicit relationship argument names.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `awx-relationship-resources`: Change relationship resource schema requirements to use explicit object-specific `_id` arguments rather than generic directional names.
- `awx-provider-documentation-and-examples`: Update generated docs/examples so relationship resource usage and breaking-change migration guidance reflect the new explicit argument names.
- `awx-typed-reference-ids`: Extend typed reference ID requirements to relationship resource reference arguments so explicit `*_id` fields remain consistently typed and directly wireable.

## Impact

- Affected code: relationship manifest derivation and naming metadata (`internal/openapi/*`, `internal/manifest/*`) and relationship resource schema/runtime mapping (`internal/provider/relationship_resource.go` and related helpers/tests).
- Affected outputs: generated manifests and documentation for relationship resources (`internal/manifest/*.json`, `docs/resources/*.md`, examples/tests that reference relationship arguments).
- API/state contracts: AWX API payload and relationship import ID formats remain unchanged; Terraform-facing argument naming becomes explicit and user-friendly.
- User impact: existing configurations that set `parent_id`/`child_id` must be updated to canonical explicit `*_id` argument names before upgrade.
