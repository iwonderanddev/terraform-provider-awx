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

## Backlog (not this change)

- [ ] Implement `awx_project.local_path` contract per [`specs/awx-project-local-path-follow-up/spec.md`](specs/awx-project-local-path-follow-up/spec.md) under the native survey initiative.
