# Tasks: Improve Provider Documentation Quality

## 1. Docs Metadata Foundation

- [ ] 1.1 Define a curated docs-enrichment data model for object docs
  (description overrides, concept-primer flag/content, further-reading links,
  and optional extra examples).
- [ ] 1.2 Add docs-enrichment loading/parsing in `cmd/awxgen` with clear errors
  for invalid metadata.
- [ ] 1.3 Add tests for docs-enrichment parsing and merge precedence behavior.

## 2. Generator Rendering Updates

- [ ] 2.1 Implement deterministic description resolution order (curated docs
  override -> manifest/OpenAPI description -> meaningful typed fallback).
- [ ] 2.2 Remove generic placeholder output
  (`Managed field from AWX OpenAPI schema`) from generated docs.
- [ ] 2.3 Update generated resource doc structure to include stable sections:
  `Example Usage`, `Schema`, `Import`, and `Further Reading`.
- [ ] 2.4 Ensure qualifier guidance is rendered separately from argument list
  entries.
- [ ] 2.5 Remove legacy directional migration phrasing and directional
  placeholder names from active relationship docs/messages.

## 3. Priority Resource Content Rollout

- [ ] 3.1 Curate field descriptions and concept/example content for
  `awx_job_template`.
- [ ] 3.2 Curate field descriptions and concept/example content for
  `awx_workflow_job_template`.
- [ ] 3.3 Curate field descriptions and concept/example content for
  `awx_project`.
- [ ] 3.4 Curate field descriptions and concept/example content for
  `awx_credential`.
- [ ] 3.5 Curate field descriptions and concept/example content for
  `awx_inventory`.
- [ ] 3.6 Curate field descriptions and concept/example content for
  `awx_inventory_source`.
- [ ] 3.7 Ensure each prioritized resource has 1-3 `Example Usage` examples and
  complex resources include concise AWX concept primers.

## 4. Docs Validation and Test Gates

- [ ] 4.1 Extend docs validation to enforce prioritized-resource rules:
  no placeholder description text, required section shape, qualifier placement,
  and example count bounds.
- [ ] 4.2 Add tests validating `Further Reading` content policy (official AWX +
  HashiCorp/AWS references where configured).
- [ ] 4.3 Add/adjust generator tests for relationship doc wording and neutral ID
  placeholder rendering.

## 5. Phase 2 All-At-Once Rollout

- [ ] 5.1 Apply the same quality gates and rendering standards to all remaining
  resource and data source docs in one batch.
- [ ] 5.2 Regenerate docs and resolve resulting validation/test regressions
  across the full documentation surface.

## 6. Final Validation and Sync

- [ ] 6.1 Run `make docs` and `make docs-validate`.
- [ ] 6.2 Run `make test`.
- [ ] 6.3 Run markdown quality checks on changed Markdown and patch any
  remaining non-fixable lint errors.
- [ ] 6.4 Confirm generated docs/manifests/spec artifacts are in sync and ready
  for `/opsx:apply`.
