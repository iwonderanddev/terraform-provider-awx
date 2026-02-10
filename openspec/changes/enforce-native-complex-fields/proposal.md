## Why

The provider currently represents many AWX `object` fields as Terraform strings. Users must wrap these values with `jsonencode(...)`, which adds unnecessary friction and obscures intent in configurations.

This is especially painful for fields such as:

- credential inputs/injectors
- notification template messages/configuration
- settings maps (for example `social_auth_google_oauth2_organization_map`)

Additionally, both `awx_job_template.extra_vars` and `awx_workflow_job_template.extra_vars` should be managed as structured values, but are currently exposed as string-based content.

## What Changes

- Treat all generated AWX `object` fields as native Terraform object values instead of JSON strings.
- Remove JSON-string compatibility for `object` fields and enforce typed values only.
- Update runtime conversion logic so provider payload/state handling for `object` fields uses native Terraform object values rather than JSON string decode/encode paths.
- Update `awx_job_template.extra_vars` and `awx_workflow_job_template.extra_vars` to be Terraform `object` fields for both:
  - `resource "awx_job_template"` and `data "awx_job_template"`
  - `resource "awx_workflow_job_template"` and `data "awx_workflow_job_template"`
- Regenerate manifests and documentation so generated resource/data source docs show the new typed behavior.
- Add/adjust provider tests for create/read/update/state behavior of object fields under the new typed-only contract.
- Keep array-field behavior out of scope for this change.

## Capabilities

### New Capabilities

- `awx-native-object-field-values`: AWX object fields are modeled as native Terraform object values rather than JSON-encoded strings.

### Modified Capabilities

- `awx-single-object-resource-model`: Clarifies that AWX object fields are configured as typed Terraform object values.
- `awx-data-sources-and-lookups`: Ensures data sources return object fields as typed Terraform object values.
- `awx-provider-documentation-and-examples`: Updates examples and argument/attribute references to remove `jsonencode(...)` patterns for managed object fields.
- `awx-sensitive-field-handling`: Preserves sensitivity/write-only handling for secret-bearing object fields under typed values.

## Impact

- Affected runtime code in `/Users/damien/git/terraform-awx-provider/internal/provider` for resource/data source schema and value conversion.
- Affected curated manifests in `/Users/damien/git/terraform-awx-provider/internal/manifest/field_overrides.json` (including job template and workflow job template `extra_vars` typing overrides).
- Affected generated outputs in `/Users/damien/git/terraform-awx-provider/internal/manifest/*.json` and `/Users/damien/git/terraform-awx-provider/docs/*`.
- Breaking behavior change for users currently passing JSON strings with `jsonencode(...)` (or raw JSON strings) for object fields.
