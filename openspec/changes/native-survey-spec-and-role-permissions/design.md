# Design: Native survey `spec` + native `role_definition.permissions`

## Context

The archived OpenSpec change `enforce-native-complex-fields` moved manifest `object` fields and job/workflow `extra_vars` to native Terraform object values. It **explicitly excluded arrays**, so `FieldTypeArray` in [`object_resource.go`](internal/provider/object_resource.go) still uses `StringAttribute` plus `decodeJSONString` / `json.Marshal` for API round-trips.

Survey-spec resources are **not** generated from the object manifest: [`relationship_resource.go`](internal/provider/relationship_resource.go) special-cases paths ending in `/survey_spec/` and hardcodes `spec` as a JSON string, so users still use `jsonencode`.

OpenAPI [`RoleDefinitionRequest`](external/awx-openapi/schema.json) defines `permissions` as an array whose items are strings—ideal for Terraform `list(string)`.

## Goals / Non-Goals

**Goals:**

- Expose survey-spec `spec` as a native structured value (same ergonomic pattern as `extra_vars`: Plugin Framework `DynamicAttribute` holding an object-shaped value; reject raw JSON strings at plan time).
- Expose `awx_role_definition.permissions` and the data source attribute as `list(string)` without `jsonencode`.
- Preserve existing import IDs and relationship identity rules (survey spec import remains parent numeric ID).
- Update provider tests and regenerated docs.

**Non-Goals:**

- General migration of **all** manifest `array` fields to native lists in this change (optional follow-up; see Decisions).
- Changing AWX API request shapes beyond encoding the same logical payloads the provider already sends.

## Decisions

1. **Survey `spec` type: `DynamicAttribute` (object)**  
   **Rationale:** Matches existing patterns for arbitrary nested objects (`extra_vars`, manifest `object` fields). [`surveySpecConfig`](internal/provider/relationship_resource.go) today decodes JSON string then POSTs decoded `any`; refactor to read `types.Dynamic`, convert with `terraformDynamicObjectToMap` (or equivalent), and POST the resulting map. Read/refresh: build Dynamic from API `map[string]any` using the same helpers as object resources.  
   **Alternative considered:** Keep string and add a second attribute—rejected (duplication and confusion).

2. **Arrays: narrow rollout for `role_definitions.permissions` only**  
   **Rationale:** Reduces breaking blast radius; OpenAPI confirms `list(string)`. Implement by branching on `(objectName == "role_definitions" && fieldName == "permissions")` in schema construction and in `payloadFromConfig` / `toTerraformValue` / `pruneUnchangedFieldsFromPayload` (and data source equivalents) to use `ListAttribute` with string element type instead of string JSON.  
   **Alternative considered:** Convert every `FieldTypeArray` field in one release—rejected unless item types are modeled for all arrays; defer to a follow-up audit of `managed_objects.json` `"type": "array"`.

3. **Breaking change, no string compatibility**  
   **Rationale:** Consistent with `enforce-native-complex-fields`. Users migrate configs in a single provider upgrade.

4. **Documentation**  
   Regenerate with `make docs`; mention migration in change notes / provider docs as needed.

## Risks / Trade-offs

- **[Risk]** Users with large existing `jsonencode` blocks see plan-time type errors. → **Mitigation:** Document migration examples (object literal for `spec`, bracket list for `permissions`).
- **[Risk]** Dynamic `spec` state drift vs. prior normalized JSON string. → **Mitigation:** Normalize from API map to a stable object representation in state; add tests for create/read/update.
- **[Risk]** Only `permissions` gets native list while other array fields remain strings. → **Mitigation:** Document clearly; track follow-up for global `FieldTypeArray` typing.

## Migration Plan

1. Upgrade provider version that includes this change.
2. Replace `spec = jsonencode({ ... })` with `spec = { ... }` (or equivalent object syntax).
3. Replace `permissions = jsonencode(["a", "b"])` with `permissions = ["a", "b"]`.
4. Run `terraform plan` and fix any remaining type diagnostics.

Rollback: pin previous provider version and revert configuration to `jsonencode` forms.

## Open Questions

- Whether to extend native list typing to all `FieldTypeArray` fields in a subsequent change once element types are catalogued from OpenAPI.
