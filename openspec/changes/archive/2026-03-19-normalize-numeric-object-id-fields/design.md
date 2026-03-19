# Design: Normalize Numeric `object_id` Fields

## Context

The provider already defines a broad contract that numeric AWX identifiers
should be exposed as Terraform numbers, with collection-created object `id`
values and canonical numeric references flowing directly between resources and
data sources. Generated role assignment surfaces currently violate that
contract: `object_id` is emitted as a string even though it represents the
numeric primary key of the object the role applies to. The mismatch is visible
in generated manifests and docs for `awx_role_team_assignment`, and the same
schema pattern may appear on related generated resources such as
`awx_role_user_assignment`.

This change is intentionally spec-driven rather than a one-off schema patch. If
the provider only special-cases one resource, the next generated surface with
the same semantic field will drift again.

## Goals / Non-Goals

**Goals:**

- Define a semantic typing rule for generated `object_id` fields that represent
  AWX numeric primary keys.
- Ensure the rule applies consistently across resource arguments, computed
  resource attributes, and data source attributes.
- Preserve nullable behavior for system-wide role assignments where `object_id`
  may be omitted.
- Document the breaking schema change and expected migration from quoted string
  literals to numeric expressions.

**Non-Goals:**

- Changing resource `id` typing or import identifier formats.
- Retyping every field literally named `object_id` without checking whether it
  semantically represents a numeric AWX key.
- Changing UUID-based alternatives such as `object_ansible_id`.

## Decisions

### Decision: Fix semantic `object_id` typing in generator-driven metadata

Implement the correction where manifest field types are derived or overridden,
so generated provider schemas, tests, and docs all inherit the same numeric
typing behavior.

Alternatives considered:

- Patch `awx_role_team_assignment` by hand in provider schema construction.
  Rejected because generated manifests and docs would remain wrong and similar
  resources would continue to drift.
- Retype every `object_id` field globally.
  Rejected because some `object_id` fields may be non-numeric in other AWX
  schemas and should remain strings.

### Decision: Scope the change to semantically numeric AWX keys

The implementation should inspect generated surfaces and only retype
`object_id` where the field represents an AWX numeric object primary key. The
design does not assume all current or future `object_id` fields are numeric.

Alternative considered:

- Restrict the fix to team assignments only.
  Rejected because the provider would keep inconsistent typing on related
  resources such as user assignments.

### Decision: Treat nullable numeric `object_id` fields as optional numbers

Fields like role-assignment `object_id` remain optional to support
system-scoped assignments, but when a value is present it must be Terraform
`Number`.

Alternative considered:

- Preserve string typing because the field may be null.
  Rejected because nullability does not require string encoding for numeric
  values.

## Risks / Trade-offs

- [Risk] The schema correction is breaking for users who currently write
  quoted `object_id` literals or rely on string-typed state.
  -> Mitigation: call out migration expectations in proposal, design, docs, and
  test coverage.

- [Risk] A broad typing rule could incorrectly retype a non-numeric
  `object_id`.
  -> Mitigation: scope the rule to semantic numeric AWX keys and add negative
  coverage for non-numeric cases if any are identified during implementation.

- [Risk] Generated docs or manifests may drift if only provider code changes.
  -> Mitigation: require the full generate/validate/docs/test/build command
  chain in the task list.

## Migration Plan

1. Identify generated resources and data sources where `object_id` represents
   an AWX numeric primary key.
2. Update generator or curated manifest typing so those fields emit numeric
   metadata.
3. Regenerate manifests and docs, then update provider schema expectations and
   tests.
4. Validate with `make generate`, `make validate-manifest`, `make docs`,
   `make docs-validate`, `make test`, and `make build`.
5. Communicate migration impact: quoted `object_id` literals must become
   numeric literals or numeric expressions, while imports and resource identity
   remain unchanged.

Rollback:

- Revert the typing-rule change and regenerate manifests/docs to restore the
  prior string schema if downstream compatibility issues outweigh the fix.

## Open Questions

- None. The change is decision-complete for implementation; the only discovery
  item is the final set of generated surfaces affected by the semantic numeric
  `object_id` rule.
