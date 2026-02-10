## Context

The provider currently maps AWX object-like fields to Terraform strings and performs JSON decode/encode internally. This requires `jsonencode(...)` in user configuration and creates unnecessary friction.

The change introduces a typed contract for all generated AWX `object` fields:

- Terraform configuration and state use native object values.
- JSON string compatibility is intentionally removed.

This is cross-cutting because it affects manifest interpretation, resource schema generation, data source schema generation, request payload conversion, state conversion, write-only preservation behavior, tests, and generated docs.

A special case is `extra_vars` on `job_templates` and `workflow_job_templates`: AWX schema models it as string, but this change requires Terraform object semantics for both resources and data sources.

## Goals / Non-Goals

**Goals:**
- Enforce native Terraform object values for all generated `object` fields in resources and data sources.
- Remove JSON-string input/output behavior for object fields.
- Preserve sensitivity and write-only handling for secret-bearing object fields.
- Type `awx_job_template.extra_vars` and `awx_workflow_job_template.extra_vars` as Terraform object fields for both resources and data sources.
- Keep generated manifests/docs/tests aligned with the new contract.

**Non-Goals:**
- No changes to generated `array` field behavior in this change.
- No backward compatibility layer for JSON-string object field inputs.
- No expansion into strongly typed nested sub-schemas per field (shape remains open object semantics).

## Decisions

### 1) Represent provider object fields using framework dynamic attributes with object-only enforcement

Decision:
- Use `DynamicAttribute` for generated `object` fields in resource and data source schemas.
- Enforce object-only values in runtime conversion (reject non-object underlying types).

Rationale:
- AWX object fields are frequently open-ended maps without fixed nested schema.
- Terraform Plugin Framework does not provide a practical static `object(any)` attribute for arbitrary nested keys/types.
- `DynamicAttribute` supports arbitrary nested values while letting runtime enforce root object shape.

Alternatives considered:
- `ObjectAttribute` with static attr map: rejected because keys/shape are not fixed.
- `MapAttribute` with a fixed element type: rejected because values are heterogeneous (nested object/list/string/number/bool/null).
- Keep string+JSON: rejected because it preserves the current usability problem.

### 2) Remove JSON decode/encode path for generated object fields

Decision:
- For generated object fields, stop reading config as `types.String` and stop serializing state as JSON strings.
- Convert object fields between Terraform dynamic/object values and Go `map[string]any` payloads.

Rationale:
- Enforces the new typed contract directly in schema/runtime.
- Eliminates string parsing ambiguity and plan-time type mismatch noise.

Alternatives considered:
- Dual mode (accept both string and object): rejected per explicit no-backward-compatibility requirement.

### 3) Preserve write-only semantics for object fields using typed values

Decision:
- Extend existing write-only snapshot preservation logic to include object fields as typed dynamic/object values.
- On read, keep prior planned/state values for write-only object fields; do not repopulate from API.

Rationale:
- Maintains existing provider invariant for secret-bearing fields.
- Avoids state loss for write-only object payloads.

Alternatives considered:
- Recompute write-only fields from API responses: rejected due to security/invariant violations.

### 4) Treat `extra_vars` as Terraform object with AWX string transport bridge

Decision:
- Add curated overrides to type:
  - `job_templates.extra_vars` as `object`
  - `workflow_job_templates.extra_vars` as `object`
- Serialize Terraform object values to JSON string when sending payloads for these fields.
- On read, accept AWX response forms and normalize to Terraform object:
  - object/map response: use directly
  - string response: parse as JSON object, with YAML fallback; empty string becomes null object field state
  - reject values that do not normalize to an object at the root

Rationale:
- Meets user-facing requirement to configure `extra_vars` as object.
- Preserves compatibility with AWX endpoint transport behavior that commonly uses string for `extra_vars`.
- Handles AWX deployments that return YAML-formatted `extra_vars` strings.

Alternatives considered:
- Keep `extra_vars` string only: rejected by requirements.
- Require AWX to return object shape: rejected because AWX contract is string-oriented here.
- JSON-only parsing without YAML fallback: rejected because it would fail valid YAML response strings in the field.

### 5) Regenerate docs from manifests and align tests with typed object behavior

Decision:
- Use existing generator/docs pipeline to propagate type changes.
- Update/add unit tests for:
  - object-field schema types
  - object-field payload conversion and validation
  - typed write-only preservation
  - `extra_vars` string bridge normalization

Rationale:
- Keeps source-of-truth model intact and avoids hand-edited docs.
- Ensures regressions are caught in conversion logic where risk is highest.

Alternatives considered:
- Manual docs edits: rejected; generated docs are authoritative in this repository.

## Risks / Trade-offs

- [Risk] Existing configurations using `jsonencode(...)` for object fields will fail type checks. → Mitigation: clearly document the breaking change and provide object-literal examples in regenerated docs.
- [Risk] YAML fallback parsing can introduce ambiguity (for example scalar/list YAML where object is required). → Mitigation: enforce object-root normalization and emit precise diagnostics for non-object results.
- [Risk] Dynamic/object conversion logic can introduce subtle plan/state diffs if normalization is inconsistent. → Mitigation: canonicalize map/list conversion paths and add targeted conversion tests.
- [Risk] Write-only object fields could accidentally be cleared during read if snapshot handling misses dynamic types. → Mitigation: extend snapshot tests to cover write-only object fields and enforce preservation behavior.

## Migration Plan

1. Add/adjust field overrides so `job_templates.extra_vars` and `workflow_job_templates.extra_vars` are `object`.
2. Update provider schema builders to emit dynamic attributes for generated `object` fields (resource + data source).
3. Implement object-only runtime conversion helpers:
   - Terraform object/dynamic -> payload `map[string]any`
   - API value -> Terraform dynamic object value
4. Add `extra_vars` transport bridge logic (object <-> AWX string).
5. Update tests for schema typing and conversion behavior.
6. Run repository validation pipeline:
   - `make generate`
   - `make validate-manifest`
   - `make docs`
   - `make docs-validate`
   - `make test`
7. Rollback strategy if needed: revert this change set to restore previous string-based object field behavior and regenerated manifests/docs.

## Open Questions

- None at this stage.
