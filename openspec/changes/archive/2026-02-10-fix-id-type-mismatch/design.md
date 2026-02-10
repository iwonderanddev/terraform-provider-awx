## Context

The provider currently has inconsistent typing between AWX/OpenAPI integer identifiers and Terraform schema types. In practice, collection-created AWX objects are integer-keyed in API payloads, while Terraform schemas historically exposed IDs and related references as strings in many places. This forces users to cast values with `tonumber(...)` when wiring resources together and makes the provider behavior diverge from the source OpenAPI contract. The change must align typing with AWX integer semantics without breaking detail-path keyed resources that intentionally use non-numeric IDs.

## Goals / Non-Goals

**Goals:**
- Define and enforce a single typing contract for AWX numeric reference fields across object resources, data sources, and relationship resources.
- Align collection-created object/data-source `id` with AWX integer semantics while preserving string IDs for detail-path keyed resources.
- Ensure generated docs reflect the new typing contract and examples no longer require explicit numeric casting for common references.
- Add tests that prevent regression in cross-resource reference typing.

**Non-Goals:**
- Redesign AWX endpoint coverage or resource naming.
- Change detail-path identifier contracts or relationship import ID formats (composite `<parent_id>:<child_id>` and survey-spec parent key).
- Introduce automatic state migration for arbitrary user-defined schema overrides beyond this typing contract.

## Decisions

### Decision: Type collection-created object/data-source `id` as Terraform number
- Rationale: For collection-created AWX objects, OpenAPI defines `id` as integer and users expect number semantics. Typing these IDs as numbers removes conversion friction and aligns the provider with AWX API contracts.
- Alternative considered: Keep all object `id` values as string. Rejected because it diverges from AWX typing and preserves counterintuitive UX.

### Decision: Keep detail-path keyed object `id` as Terraform string
- Rationale: Some managed objects are keyed by non-numeric path identifiers (for example settings categories), so forcing numeric type would be incorrect.
- Alternative considered: Type every object `id` as number. Rejected because detail-key endpoints are not integer keyed.

### Decision: Treat AWX integer foreign-key/reference fields as Terraform numbers in both arguments and attributes
- Rationale: Users should be able to pass values directly between resources/data sources without manual conversion. This aligns plan-time types and reduces configuration noise.
- Alternative considered: Keep current mixed typing and document `tonumber(...)` usage. Rejected because this preserves avoidable friction and weakens provider UX.

### Decision: Centralize type normalization in manifest/openapi derivation, not ad hoc per-resource patches
- Rationale: Provider surfaces are generated dynamically from manifests. A single normalization step ensures consistency across all generated resources/data sources and avoids drift.
- Alternative considered: Patch individual resources with overrides. Rejected because it is brittle and does not scale with schema refreshes.

### Decision: Add compatibility guardrails for fields already persisted as strings
- Rationale: Some existing states may have string-typed values for reference fields and object IDs. Implementation should avoid destructive behavior by preferring framework-compatible conversions where possible and documenting any unavoidable upgrade behavior.
- Alternative considered: Hard switch without compatibility handling. Rejected due to upgrade friction risk.

## Risks / Trade-offs

- [State/type transition for existing users] -> Mitigation: add upgrade-focused tests for collection-created object/data-source `id` transitions and document expected plan behavior.
- [Incorrect inference of numeric vs non-numeric fields from schema metadata] -> Mitigation: constrain normalization to known integer field shapes and validate via manifest checks/tests.
- [Generated docs lag behind schema changes] -> Mitigation: keep `make generate`, `make docs`, and validation steps mandatory in implementation tasks.
- [Relationship resources with composite string IDs may be conflated with numeric references] -> Mitigation: explicitly exclude relationship identity/composite import IDs from numeric normalization.
