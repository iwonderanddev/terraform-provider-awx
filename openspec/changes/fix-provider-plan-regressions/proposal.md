## Why

`terraform plan` failed in `/Users/damien/git/mockshop-in-cloud-2/terraform/infrastructure/awx/dev` due to two provider regressions:

- Type conversion failures in `data.awx_credential.galaxy_default` for write-only integer fields (`team`, `user`) being written as string nulls.
- Name lookup pagination failures in `data.awx_role_definition` when AWX returned query-only `next` links (for example `?page=2`), which were incorrectly encoded into the request path.

These are provider correctness issues and must be fixed to restore deterministic planning behavior.

## What Changes

- Fix write-only field state handling to be type-aware across object resources and data sources.
- Preserve write-only planned/state values using native Terraform value types instead of assuming string-only handling.
- Normalize pagination `next` URL resolution so absolute, relative, and query-only forms are all handled correctly.
- Add regression tests for:
  - Write-only integer state null/default handling.
  - Mixed-type write-only value preservation.
  - Query-only AWX pagination links in `ListAll`.

## Impact

- Removes provider-side `Value Conversion Error` failures for `awx_credential` data source reads.
- Restores reliable multi-page lookup behavior for objects such as `role_definitions`.
- Improves type safety for future write-only fields across generated object surfaces.
