# Design: Rename Project Credential ID To SCM Credential ID

## Context

The AWX `projects` object exposes a numeric `credential` field that AWX uses
for source control access. The provider currently applies its generic
reference-field naming rule and emits that field as `credential_id` on both the
`awx_project` resource and data source.

That generated name is technically consistent with the existing `_id` suffix
rule, but it is not semantically precise. Within the provider surface,
credential references that have distinct purposes already use purpose-specific
names such as `webhook_credential_id`,
`signature_validation_credential_id`, and `source_credential_id`. The project
SCM credential is the outlier.

Because the user requested a breaking rename, this design assumes there is no
compatibility alias. The canonical Terraform field must change from
`credential_id` to `scm_credential_id` wherever the generated `awx_project`
schema is emitted.

## Goals / Non-Goals

**Goals:**

- Make the `awx_project` Terraform schema explicitly communicate that the field
  is the SCM credential.
- Apply the rename consistently to both the resource and data source surfaces.
- Encode the rename in generator metadata so the field name is derived
  declaratively and remains visible in generated manifests.
- Regenerate docs and examples so they teach the renamed field consistently.
- Preserve the underlying AWX API contract and state/import identity behavior.

**Non-Goals:**

- Introducing a compatibility alias or deprecation window for
  `awx_project.credential_id`.
- Renaming other generic `credential_id` fields across unrelated resources.
- Changing AWX API payload shape or the provider's CRUD/import semantics for
  projects.
- Redesigning reference-field naming globally beyond the metadata needed for
  this class of override.

## Decisions

### Decision: Add declarative canonical field-name overrides for object fields

Introduce generator metadata that allows a managed object field to publish a
canonical Terraform name different from the raw AWX field-derived default.

For this change, the AWX project field `credential` will declare
`scm_credential_id` as its canonical Terraform name.

Rationale:

- Keeps the naming decision close to curated metadata rather than embedding a
  one-off `if project && credential` rule in runtime code.
- Scales to future cases where AWX field names are technically references but
  semantically ambiguous in Terraform.
- Ensures generated manifests remain the source of truth for downstream schema
  generation and docs.

Alternatives considered:

- Hard-code the rename in provider schema assembly.
  Rejected because it hides naming policy outside the generation pipeline.
- Only change docs text while keeping `credential_id`.
  Rejected because the user wants Terraform configuration itself to reflect SCM
  intent.

### Decision: Treat `scm_credential_id` as the only supported Terraform field

The provider will remove `credential_id` from the generated `awx_project`
resource and data source schemas rather than keeping both names in parallel.

Rationale:

- Matches the requested breaking rename directly.
- Avoids ambiguous dual-field validation and plan/state precedence rules.
- Keeps the canonical field contract simple for generated docs and tests.

Alternatives considered:

- Keep `credential_id` as a deprecated alias during transition.
  Rejected because it postpones the rename without eliminating ambiguity.

### Decision: Update documentation curation metadata alongside manifest output

The `awx_project` docs-enrichment entry and any generated examples must change
to `scm_credential_id` in the same change as the manifest/schema rename.

Rationale:

- Prevents generated docs from lagging behind the renamed schema.
- Keeps examples immediately runnable against the new contract.

Alternatives considered:

- Rely on generated manifests alone and patch docs later.
  Rejected because it creates an inconsistent proposal and incomplete task
  breakdown.

## Risks / Trade-offs

- [Risk] Adding canonical-name override metadata could be broader than this one
  rename and invite under-specified future overrides.
  -> Mitigation: scope the override model to explicit per-object field entries
  and cover its behavior with targeted generator tests.

- [Risk] The breaking rename will invalidate existing Terraform configuration
  and state expectations for `awx_project`.
  -> Mitigation: document the rename explicitly in proposal/spec/tasks and
  update examples and acceptance coverage to make the new field obvious.

- [Risk] Generated outputs may still contain `credential_id` in curated
  documentation metadata or tests even after schema derivation changes.
  -> Mitigation: include docs-enrichment, generated docs, and schema assertions
  in the implementation task list.

## Migration Plan

1. Add object-field canonical-name override support to the generation pipeline.
2. Declare an override for `projects.credential` so generated Terraform field
   naming becomes `scm_credential_id`.
3. Regenerate manifests and confirm `awx_project` resource/data source metadata
   no longer exposes `credential_id`.
4. Update curated docs metadata and regenerate docs/examples.
5. Add or update tests that assert the renamed field is present and the legacy
   name is absent.
6. Run the documented generate/validate/docs/test/build command chain.

Rollback:

- Revert the override support and the `projects.credential` rename metadata,
  then regenerate manifests and docs to restore `credential_id`.

## Open Questions

None.
