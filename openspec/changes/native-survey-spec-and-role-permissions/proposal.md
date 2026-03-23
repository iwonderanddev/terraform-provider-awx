# Proposal: Native survey `spec` + native `role_definition.permissions`

## Why

Survey-spec relationship resources (`awx_job_template_survey_spec`, `awx_workflow_job_template_survey_spec`) expose `spec` as a Terraform string, so users must wrap payloads with `jsonencode(...)`. The archived change `enforce-native-complex-fields` removed that pattern for manifest `object` fields and `extra_vars`, but survey specs are implemented outside the manifest object path in `relationship_resource.go`, so they were never updated.

Separately, `awx_role_definition.permissions` is a manifest `array` of strings in the OpenAPI schema, but the provider still maps all `FieldTypeArray` fields through JSON-string transport in `object_resource.go` (the prior change explicitly kept arrays out of scope). Users must `jsonencode` list payloads instead of using native `list(string)`.

## What Changes

- **BREAKING**: Change survey-spec `spec` from `String` (JSON) to a native structured value (aligned with `extra_vars`: Terraform Plugin Framework `Dynamic` / object semantics), for both job-template and workflow job-template survey-spec resources.
- **BREAKING**: Change `awx_role_definition.permissions` (and the corresponding data source) from JSON string to `list(string)` per OpenAPI `RoleDefinitionRequest.permissions` (`items: string`).
- Document migration: replace `spec = jsonencode({...})` with object literals; replace `permissions = jsonencode([...])` with `permissions = [...]`.
- Implementation scope for arrays: **narrow** first—native `list(string)` for `role_definitions.permissions` only unless implementation generalizes all `FieldTypeArray` string fields in the same change (design.md records the chosen approach).

## Capabilities

### New Capabilities

- `awx-native-array-field-values`: Manifest `array` fields with string element types SHALL be exposed as native Terraform lists (starting with `role_definitions.permissions`), not JSON-encoded strings.

### Modified Capabilities

- `awx-relationship-resources`: Survey-spec `spec` SHALL be a structured Terraform value, not a JSON string.
- `awx-native-object-field-values`: Survey-spec relationship `spec` payloads SHALL follow native structured object semantics consistent with other non-string complex values.
- `awx-single-object-resource-model`: Clarify that manifest `array` fields use native list/tuple semantics where specified, not JSON-string configuration for those fields.
- `awx-provider-documentation-and-examples`: Examples and argument descriptions SHALL not require `jsonencode` for survey `spec` or `awx_role_definition.permissions`.

## Impact

- Runtime: [`internal/provider/relationship_resource.go`](internal/provider/relationship_resource.go) (survey-spec schema, create/read/update/state).
- Runtime: [`internal/provider/object_resource.go`](internal/provider/object_resource.go), [`internal/provider/object_data_source.go`](internal/provider/object_data_source.go) (array field schema and value conversion for `permissions` and any shared `FieldTypeArray` paths).
- Tests: [`internal/provider/relationship_resource_test.go`](internal/provider/relationship_resource_test.go), [`internal/provider/acceptance_terraform_test.go`](internal/provider/acceptance_terraform_test.go), array/object conversion tests as needed.
- Regenerated [`docs/*`](docs/) after `make docs`.
- Breaking for existing modules using string/`jsonencode` for these attributes.
