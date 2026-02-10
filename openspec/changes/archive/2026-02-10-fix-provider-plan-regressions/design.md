## Context

The provider currently generates object resources and data sources from manifest metadata. Field typing is schema-driven (`integer`, `boolean`, `number`, `string`/JSON-encoded complex types). Write-only fields exist across several objects, including integer write-only fields on `credentials` (`team`, `user`).

The AWX API client uses `ListAll` to traverse paginated collection responses by following AWX `next` links.

## Goals / Non-Goals

**Goals:**

- Keep write-only field state operations type-correct for all field kinds.
- Ensure pagination works when AWX returns absolute URLs, relative paths, or query-only links.
- Add focused regression coverage for both failure classes.

**Non-Goals:**

- No change to import ID contracts.
- No change to manifest generation behavior.
- No authentication or transport feature expansion.

## Decisions

### 1) Write-only state storage must be typed

Decision:
- Store write-only snapshots as typed Terraform values (`types.Int64`, `types.Bool`, `types.Float64`, `types.String`) rather than string-only values.
- When no write-only value is available, write a typed null based on field metadata.

Rationale:
- Prevents framework type mismatch errors when the schema expects numeric/boolean values.
- Preserves current secret/write-only behavior while supporting non-string write-only fields.

### 2) Data source write-only fields use typed null conversion

Decision:
- For write-only attributes in object data sources, set null state values via shared field-aware conversion logic.

Rationale:
- Keeps data source behavior aligned with resource state behavior.
- Avoids ad hoc null assignment that assumes a string schema type.

### 3) Pagination `next` handling resolves via URL reference semantics

Decision:
- Resolve `next` against current request path/query using URL reference resolution instead of raw string path substitution.

Rationale:
- AWX pagination can return `next` as:
  - Absolute URL (`https://.../api/v2/x/?page=2`)
  - Relative path (`/api/v2/x/?page=2`)
  - Query-only (`?page=2`)
- Reference resolution handles all forms safely and prevents malformed requests such as `/api/v2/x/%3Fpage=2`.

### 4) Regression tests are required for both bugs

Decision:
- Add tests for typed write-only state handling and query-only pagination.

Rationale:
- Both issues were discovered by live Terraform plan execution and should remain guarded in unit tests.

## Risks / Trade-offs

- [Type conversion coupling] Write-only logic now branches by field type. Mitigation: reuse existing conversion/null helpers and add targeted tests.
- [Pagination edge cases] URL resolution could alter behavior for unusual `next` formats. Mitigation: keep compatibility tests for existing absolute URL behavior and add query-only coverage.
