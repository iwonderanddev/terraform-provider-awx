# Tasks: Improve Provider Documentation Quality

## 1. Docs Metadata Foundation

- [x] 1.1 Define a curated docs-enrichment data model for object docs
  (description overrides, concept-primer flag/content, further-reading links,
  and optional extra examples).
- [x] 1.2 Add docs-enrichment loading/parsing in `cmd/awxgen` with clear errors
  for invalid metadata.
- [x] 1.3 Add tests for docs-enrichment parsing and merge precedence behavior.

## 2. Generator Rendering Updates

- [x] 2.1 Implement deterministic description resolution order (curated docs
  override -> manifest/OpenAPI description -> meaningful typed fallback).
- [x] 2.2 Remove generic placeholder output
  (`Managed field from AWX OpenAPI schema`) from generated docs.
- [x] 2.3 Update generated resource doc structure to include stable sections:
  `Example Usage`, `Schema`, `Import`, and `Further Reading`.
- [x] 2.4 Ensure qualifier guidance is rendered separately from argument list
  entries.
- [x] 2.5 Remove legacy directional migration phrasing and directional
  placeholder names from active relationship docs/messages.

## 3. Priority Resource Content Rollout

- [x] 3.1 Curate field descriptions and concept/example content for
  `awx_job_template`.
- [x] 3.2 Curate field descriptions and concept/example content for
  `awx_workflow_job_template`.
- [x] 3.3 Curate field descriptions and concept/example content for
  `awx_project`.
- [x] 3.4 Curate field descriptions and concept/example content for
  `awx_credential`.
- [x] 3.5 Curate field descriptions and concept/example content for
  `awx_inventory`.
- [x] 3.6 Curate field descriptions and concept/example content for
  `awx_inventory_source`.
- [x] 3.7 Ensure each prioritized resource has 1-3 `Example Usage` examples and
  complex resources include concise AWX concept primers.

## 4. Docs Validation and Test Gates

- [x] 4.1 Extend docs validation to enforce prioritized-resource rules:
  no placeholder description text, required section shape, qualifier placement,
  and example count bounds.
- [x] 4.2 Add tests validating `Further Reading` content policy (official AWX
  links only, including resource specificity and non-generic index handling).
- [x] 4.3 Add/adjust generator tests for relationship doc wording and neutral ID
  placeholder rendering.
- [x] 4.4 Add docs-enrichment provenance validation for prioritized resources
  (`officialAwxUrl` + `verifiedOn`) and corresponding tests.

## 5. Phase 2 All-At-Once Rollout

- [x] 5.1 Apply the same quality gates and rendering standards to all remaining
  resource and data source docs in one batch.
- [x] 5.2 Regenerate docs and resolve resulting validation/test regressions
  across the full documentation surface.

## 6. Final Validation and Sync

- [x] 6.1 Run `make docs` and `make docs-validate`.
- [x] 6.2 Run `make test`.
- [x] 6.3 Run markdown quality checks on changed Markdown and patch any
  remaining non-fixable lint errors.
- [x] 6.4 Confirm generated docs/manifests/spec artifacts are in sync and ready
  for `/opsx:apply`.

## 7. Official AWX Doc Grounding and Link Specificity

- [x] 7.1 Build and verify an official AWX 24.6.1 per-object documentation link
  mapping from live online docs.
- [x] 7.2 Update docs generation so each managed object resource/data source
  includes a resource-specific official AWX link in `## Further Reading` (not
  only generic index links).
- [x] 7.3 Update relationship docs `## Further Reading` to point to mapped
  parent/child official AWX concept links where available.
- [x] 7.4 Extend docs validation tests/rules to enforce AWX link specificity and
  reject generic-index-only AWX linking.
- [x] 7.5 Re-curate prioritized resource examples/descriptions where official
  AWX 24.6.1 pages indicate wording or behavior adjustments.
- [x] 7.6 Regenerate docs and rerun `make docs`, `make docs-validate`,
  `make test`, and markdown lint on changed Markdown.
