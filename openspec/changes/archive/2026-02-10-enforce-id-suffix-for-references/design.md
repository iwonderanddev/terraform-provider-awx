## Context

The provider currently exposes many AWX foreign-key link fields using unsuffixed names (for example, `organization`), even though the value is an object identifier. This obscures intent in Terraform configuration and makes link wiring less self-documenting than explicit names like `organization_id`.

This change is cross-cutting because field names are derived in generation and consumed in multiple places: manifest generation, provider runtime schema construction, state mapping, import/state tests, and generated docs/examples.

Constraints that must remain intact:
- Preserve object identity/import contracts (`id` semantics for collection-created vs detail-path keyed objects).
- Preserve relationship resource contracts (`parent_id`, `child_id`, survey spec handling).
- Keep generated manifests and docs as source-of-truth outputs from the existing generation pipeline.

## Goals / Non-Goals

**Goals:**
- Enforce a provider-wide naming contract: resource/data-source fields that represent AWX object links use `_id` suffixes.
- Apply naming consistently for both configurable arguments and computed attributes for those link fields.
- Ensure generated docs/examples only show canonical `_id` names for link wiring.
- Define migration behavior for users currently using unsuffixed names.

**Non-Goals:**
- No change to relationship resource argument names (`parent_id`, `child_id`) because they already express identifier intent.
- No change to non-reference fields, even if their AWX names happen to resemble IDs.
- No change to AWX API payload keys; only Terraform-facing schema names are adjusted.

## Decisions

### 1) Canonical Terraform field names for object links SHALL be suffixed with `_id`

Decision:
- For generated object resources/data sources, any field classified as an AWX object-link reference is exposed as `<name>_id`.
- Unsuffixed link field names are no longer canonical and are removed from generated schemas.

Rationale:
- `_id` fields are immediately recognizable as references, improving readability and reducing configuration ambiguity.
- A single naming rule is easier to document, validate, and test across all generated surfaces.

Alternatives considered:
- Keep unsuffixed names and rely on docs: rejected because ambiguity persists in real configs.
- Support both unsuffixed and suffixed names indefinitely: rejected due to long-term schema duplication and drift risk.

### 2) Reference-field classification SHALL come from manifest/openapi-derived metadata, not name heuristics

Decision:
- Continue using generated field metadata (relation/reference semantics) to determine whether a field is a link.
- Do not infer link-ness from string patterns like `_id` suffixes alone.

Rationale:
- Avoids false positives/negatives where AWX naming is inconsistent.
- Keeps behavior tied to source-of-truth schema and curated overrides.

Alternatives considered:
- Regex/name-based detection only: rejected as brittle and likely to regress edge endpoints.

### 3) Migration strategy is a breaking rename with explicit diagnostics and docs

Decision:
- Treat unsuffixed link fields as removed from provider schemas.
- Surface clear diagnostics when configuration still uses removed names.
- Regenerated documentation/examples become the primary migration guide.

Rationale:
- Keeps implementation simple and deterministic in generated code.
- Prevents hidden precedence/alias conflicts between old and new names.

Alternatives considered:
- One-release alias/deprecation period: rejected for now because it adds complexity across generation/runtime and weakens naming consistency guarantees.

### 4) Keep typed-ID semantics unchanged while applying them to canonical suffixed fields

Decision:
- Existing ID type rules remain: numeric references stay Terraform numbers, collection-created object `id` remains numeric, detail-path keyed object `id` remains string.
- The rename changes only field names, not underlying identifier type semantics.

Rationale:
- Preserves prior type-contract work and avoids introducing unnecessary behavioral deltas.

Alternatives considered:
- Revisit all ID typing simultaneously: rejected as scope expansion unrelated to naming objective.

## Risks / Trade-offs

- [Risk] Existing Terraform configurations using unsuffixed fields will fail plan/apply after upgrade. -> Mitigation: document breaking rename clearly and provide before/after examples for common resources.
- [Risk] Incorrect reference-field classification could rename non-reference fields or miss true references. -> Mitigation: add generator/provider tests for representative object sets and edge cases.
- [Risk] Name collisions where both `<name>` and `<name>_id` exist in source metadata could create invalid schemas. -> Mitigation: add manifest validation that detects collisions and fails generation with actionable diagnostics.
- [Risk] Docs can drift from generated schema expectations. -> Mitigation: require docs regeneration and `make docs-validate` in the delivery checklist.

## Migration Plan

1. Update manifest/openapi-to-catalog mapping so Terraform-facing names for reference fields are generated as `_id` suffixed names.
2. Regenerate manifests and ensure runtime schema builders consume canonical suffixed names.
3. Update docs/examples to remove unsuffixed link usage and show `_id` wiring.
4. Add or update provider/generator tests to enforce naming and typing invariants.
5. Run validation chain: `make generate`, `make validate-manifest`, `make docs`, `make docs-validate`, `make test`.
6. Rollback strategy: revert this change set and regenerate manifests/docs to restore previous unsuffixed schema names.

## Open Questions

- None. This change is confirmed as a strict breaking rename across all generated object-link fields.
