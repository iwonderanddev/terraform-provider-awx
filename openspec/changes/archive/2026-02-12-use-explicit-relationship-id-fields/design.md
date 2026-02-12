## Context

Relationship resources are currently generated with generic arguments (`parent_id`, `child_id`) regardless of the concrete AWX objects being linked. This makes Terraform configuration harder to read and increases miswiring risk because users must infer directionality from endpoint behavior instead of schema names.

This change affects generator, manifest, provider runtime, tests, and generated docs because relationship argument names are created in generation and consumed across the provider surface.

Constraints that must remain intact:
- Preserve relationship resource import/state identity contracts (`<parent_id>:<child_id>` for association resources and `<parent_id>` for survey-spec/singleton resources).
- Preserve AWX API transport semantics (`id` in POST/associate/disassociate payloads and parent path keying).
- Keep resource type names unchanged; this change is argument naming only.

## Goals / Non-Goals

**Goals:**
- Expose canonical relationship arguments as explicit object-specific `_id` names (for example, `job_template_id`, `credential_id`).
- Keep canonical relationship arguments consistently typed as Terraform numbers where AWX expects numeric object IDs.
- Maintain existing import formats and runtime behavior while improving schema clarity.
- Provide explicit breaking-change migration guidance from `parent_id`/`child_id` to canonical fields.

**Non-Goals:**
- No relationship resource type rename.
- No changes to object resource/data source field naming beyond what is needed for relationship-resource consistency.
- No change to survey spec `spec` JSON contract or singleton import ID format.

## Decisions

### 1) Canonical relationship arguments SHALL be object-specific `_id` fields

Decision:
- Standard relationship resources expose `<parent_object_singular>_id` and `<child_object_singular>_id` as canonical required inputs.
- Survey spec relationship resources expose `<parent_object_singular>_id` plus `spec`.

Rationale:
- The argument names become self-describing and remove directional ambiguity for users.
- Names can be generated deterministically from existing relationship manifest metadata (`parentObject`, `childObject`).

Alternatives considered:
- Keep `parent_id`/`child_id` and improve docs only: rejected because schema remains ambiguous in real configs.
- Introduce free-form aliases only in docs/examples: rejected because usability issue remains in the API contract.

### 2) Backward compatibility SHALL be a hard break with canonical fields only

Decision:
- `parent_id` and `child_id` are removed from relationship resource schemas.
- Relationship resources accept only canonical object-specific `*_id` arguments.
- Documentation and release notes include a breaking-change mapping from legacy names to canonical names.

Rationale:
- Removes ambiguous directional naming immediately and enforces a single canonical contract.
- Avoids long-term runtime complexity caused by alias branching and precedence handling.

Alternatives considered:
- One-release deprecation alias window: rejected per product decision for immediate clarity.
- Permanent dual-schema support: rejected due to long-term complexity and unclear canonical usage.

### 3) Relationship identity/import contracts SHALL remain unchanged

Decision:
- State `id` remains `<parent_id>:<child_id>` for standard associations.
- State/import `id` remains `<parent_id>` for singleton survey-spec relationships.
- Canonical argument renaming affects Terraform schema only, not ID encoding.

Rationale:
- Preserves existing state and import behavior, minimizing operational disruption.
- Separates usability improvements from identity semantics.

Alternatives considered:
- Encode named keys in state ID (for example, `job_template_id:credential_id`): rejected as unnecessary contract churn.

### 4) Generation-owned naming metadata SHALL drive runtime schema and docs

Decision:
- Canonical relationship argument names are generated and embedded in manifest data.
- Provider runtime uses generated canonical names to read config, set state, and map to AWX requests.
- Docs/examples are regenerated from the same metadata to prevent drift.

Rationale:
- Maintains a single source of truth across generator, runtime, and documentation.
- Prevents one-off hand patches in provider runtime.

Alternatives considered:
- Hard-code relationship argument names in runtime code: rejected because it duplicates generator knowledge and scales poorly.

## Risks / Trade-offs

- [Risk] Existing configurations using `parent_id`/`child_id` will fail after upgrade. -> Mitigation: provide explicit migration guidance and before/after examples in docs and changelog.
- [Risk] Generated singular names may be awkward for certain collections/endpoints. -> Mitigation: reuse existing singularization rules and allow curated overrides where needed.
- [Risk] Migration may be incomplete in acceptance tests/examples. -> Mitigation: update all relationship examples/tests and run full validation chain.

## Migration Plan

1. Extend relationship manifest generation to emit canonical explicit argument names for parent/child references.
2. Update relationship resource schema/runtime to use canonical names only.
3. Regenerate manifests and docs, and update examples/acceptance tests to canonical names.
4. Validate with `make generate`, `make validate-manifest`, `make docs`, `make docs-validate`, and `make test`.

Rollback strategy:
- Revert the change set and regenerate manifests/docs to restore `parent_id`/`child_id` as canonical arguments.

## Open Questions

- None.
