# Tasks: Native survey `spec` + `role_definition.permissions`

## 1. Survey-spec relationship resources

- [ ] 1.1 Replace `spec` `StringAttribute` with `DynamicAttribute` in [`internal/provider/relationship_resource.go`](internal/provider/relationship_resource.go) for `isSurveySpecRelationship()` schemas.
- [ ] 1.2 Refactor `surveySpecConfig`, `setSurveySpecState`, `Read`, `Create`, and `Update` to use `types.Dynamic` and existing object conversion helpers (`terraformDynamicObjectToMap`, API map to dynamic/object for state).
- [ ] 1.3 Ensure POST payload matches prior semantic behavior (decoded object to AWX JSON body).
- [ ] 1.4 Update unit tests in [`internal/provider/relationship_resource_test.go`](internal/provider/relationship_resource_test.go) and acceptance HCL in [`internal/provider/acceptance_terraform_test.go`](internal/provider/acceptance_terraform_test.go) to remove `jsonencode` for `spec`.
- [ ] 1.5 Adjust [`internal/provider/legacy_argument_diagnostics_test.go`](internal/provider/legacy_argument_diagnostics_test.go) if it asserts legacy string typing for survey `spec`.

## 2. `awx_role_definition.permissions` (narrow: `role_definitions` only)

- [ ] 2.1 In [`internal/provider/object_resource.go`](internal/provider/object_resource.go), emit `ListAttribute` with string elements for `role_definitions` + `permissions` (resource schema); keep other `FieldTypeArray` fields on string JSON until a follow-up.
- [ ] 2.2 In [`internal/provider/object_data_source.go`](internal/provider/object_data_source.go), mirror list typing for the same field on the data source.
- [ ] 2.3 Update `payloadFromConfig`, `pruneUnchangedFieldsFromPayload`, and `toTerraformValue` (and data source setState paths) for list read/write and API `[]any` / `[]string` conversion.
- [ ] 2.4 Add or extend tests (e.g. [`internal/provider/object_field_names_test.go`](internal/provider/object_field_names_test.go), payload tests) for `permissions` list schema and round-trip.

## 3. Generated docs and validation

- [ ] 3.1 Run `make generate`, `make validate-manifest`, `make docs`, `make docs-validate` per [`AGENTS.md`](AGENTS.md).
- [ ] 3.2 Run `make test` and `make build`.
- [ ] 3.3 Run `openspec validate --changes native-survey-spec-and-role-permissions` before merge.

## 4. Migration notes

- [ ] 4.1 Document breaking migration in release notes or provider changelog entry: `spec` object syntax; `permissions = [...]` instead of `jsonencode`.
