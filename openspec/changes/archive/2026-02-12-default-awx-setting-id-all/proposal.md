# Proposal: Default awx_setting ID to all in Docs/UX

## Why

The current `awx_setting` documentation and examples do not make the intended
operator workflow clear, and they increase user friction by forcing users to
pick a settings subset before they can manage configuration. We should optimize
for the common UX path by making `id = "all"` the documented default now, while
still preserving category-scoped IDs for optional advanced use.

## What Changes

- Update generated and curated `awx_setting` docs and examples to use
  `id = "all"` as the default and recommended value.
- Clarify that category IDs (for example `system`, `authentication`, `bulk`)
  remain supported as optional scoping mechanisms.
- Add explicit guidance about overlap/conflict risk when multiple
  `awx_setting` resources manage intersecting keys (especially mixing `all`
  with category-scoped resources).
- Align import guidance and examples with the default `all` UX so users can
  follow a single canonical path without additional endpoint discovery.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `awx-provider-documentation-and-examples`: Update requirement-level
  documentation expectations for `awx_setting` so default examples/import
  guidance use `id = "all"`, while category-scoped IDs are documented as
  optional advanced usage with clear overlap/conflict warnings.

## Impact

- Affected outputs: `/Users/damien/git/terraform-provider-awx-iwd/docs/resources/awx_setting.md` and
  `/Users/damien/git/terraform-provider-awx-iwd/docs/data-sources/awx_setting.md`.
- Likely affected generators/metadata: docs generation behavior in
  `/Users/damien/git/terraform-provider-awx-iwd/cmd/awxgen` and any related curation metadata in
  `/Users/damien/git/terraform-provider-awx-iwd/internal/manifest` used for examples/import messaging.
- No AWX API contract changes and no import ID format changes are required; the
  change is documentation/UX guidance and consistency.
