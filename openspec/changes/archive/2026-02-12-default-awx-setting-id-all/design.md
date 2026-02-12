# Design: Default awx_setting ID to all in Docs/UX

## Context

The provider currently supports detail-path keyed `awx_setting` resources where
`id` maps to AWX settings category slugs. In practice, AWX supports `all` as a
valid category slug that exposes the full settings payload and permits updates
across setting groups. Existing documentation does not make `all` the canonical
path, which adds unnecessary user decision overhead and creates inconsistent
examples/import guidance.

This change is documentation and generation-policy scoped. It must preserve
runtime compatibility, import ID contracts, and category-scoped usage while
standardizing docs around `id = "all"` for the default UX.

## Goals / Non-Goals

**Goals:**

- Make `id = "all"` the default and recommended `awx_setting` path in generated
  docs and examples.
- Preserve support for category-scoped IDs and document them as optional,
  advanced scoping choices.
- Add explicit warnings about potential overlap/conflicts when users manage
  intersecting settings keys across multiple `awx_setting` resources.
- Keep docs generation deterministic and validation-enforced to prevent
  regressions.

**Non-Goals:**

- No behavior changes to provider runtime CRUD logic for settings resources.
- No AWX API changes, endpoint filtering, or permission model changes.
- No deprecation/removal of category-scoped IDs.

## Decisions

1. Default UX decision: documentation defaults to `id = "all"`.
Rationale: this minimizes operator cognitive load and aligns with the most
straightforward workflow when managing settings declaratively.
Alternative considered: keep category-first guidance. Rejected because it
forces users to pre-classify keys into AWX subsets before they can start.

2. Compatibility decision: keep category-scoped IDs documented and supported.
Rationale: some users prefer tighter blast radius and ownership boundaries.
Alternative considered: hide category IDs entirely. Rejected because it removes
useful scoping patterns and existing operator workflows.

3. Safety decision: document overlap/conflict semantics explicitly.
Rationale: when `all` and scoped resources touch the same keys, Terraform state
ownership can become ambiguous. Clear warnings reduce avoidable drift and
surprises.
Alternative considered: no explicit warning. Rejected because conflict behavior
is not obvious to users.

4. Enforcement decision: update docs validation/tests for `awx_setting`
default example/import expectations.
Rationale: codifies this UX choice in CI and prevents regression to ambiguous or
placeholder setting IDs.
Alternative considered: policy by convention only. Rejected due to high
regression risk in generated docs.

## Risks / Trade-offs

- [Users interpret default as exclusive] -> Mitigation: explicitly document that
  category-scoped IDs remain supported for advanced scoping.
- [Overlap conflicts when mixing `all` and scoped resources] -> Mitigation: add
  clear warnings and recommend single-owner patterns for shared keys.
- [Docs drift from generator behavior] -> Mitigation: encode expectations in
  docs generation tests/validation.

## Migration Plan

- Regenerate docs with updated default examples/import guidance for
  `awx_setting`.
- Validate docs/test gates pass and confirm both `all` and category-scoped
  examples are accurately represented.
- Rollback strategy: revert docs-generation policy changes if validation reveals
  incompatibilities; runtime behavior remains unchanged.

## Open Questions

- Should the data source `awx_setting` also default examples to `id = "all"` in
  every case, or should some examples demonstrate scoped lookup first?
- Should warnings include a concrete “do not mix ownership for the same key”
  pattern example in docs metadata for stronger guidance?
