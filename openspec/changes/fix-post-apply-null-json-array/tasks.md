# Tasks: fix-post-apply-null-json-array

## Implementation

- [x] Add `jsonEncodedArrayStringValuesFromSource` and thread prior JSON-encoded array strings through `Create`, `Read`, `Update` into `setState`.
- [x] Implement `normalizeOptionalEmptyJSONEncodedArrayToNull` + `isEmptyJSONArrayAPIValue` for optional non-computed JSON-array fields.
- [x] Unit tests: `TestNormalizeOptionalEmptyJSONEncodedArrayToNull`, `TestSetStateOptionalJSONEncodedArrayNullWhenEmptyAPIAndPriorNull`, `TestSetStateOptionalJSONEncodedArrayKeepsExplicitEmptyJSONArray`.
- [x] Update existing `setState` call sites in tests.

## OpenSpec

- [x] Add `proposal.md`, `design.md`, `tasks.md`, `.openspec.yaml`.
- [x] Add spec deltas under `specs/`.

## Verification

- [x] `go test ./...`
- [x] `make generate` / `make validate-manifest`
- [x] `make docs` / `make docs-validate`

## `local_path` read-only (follow-up contract)

- [x] Add `FieldSpec.readOnly` and override `projects.local_path` with `computed` + `readOnly`.
- [x] Schema: `newResourceFieldAttribute` uses read-only for Optional=false, Computed=true.
- [x] Skip read-only fields in `payloadFromConfig`; `setState` preserves API value for read-only strings.
- [x] Resource docs: list read-only attributes under **Read-Only** (not Optional).
- [x] Update OpenSpec `design.md` / `proposal.md` / [`specs/awx-project-local-path-follow-up/spec.md`](specs/awx-project-local-path-follow-up/spec.md).
