# awx-native-array-field-values Delta

## ADDED Requirements

### Requirement: String array manifest fields use Terraform lists

The provider SHALL expose `role_definitions.permissions` as Terraform `list(string)` on `awx_role_definition` and `data.awx_role_definition` because OpenAPI defines `permissions` as an array whose items are strings. Configuration for this field MUST NOT accept JSON-encoded strings.

#### Scenario: Permissions list configuration

- **WHEN** a user sets `permissions` on `awx_role_definition`
- **THEN** Terraform accepts a list of string elements and SHALL reject JSON-string payloads with a type diagnostic

#### Scenario: Data source readback for permissions

- **WHEN** a data source reads `awx_role_definition` and AWX returns `permissions` as a JSON array of strings
- **THEN** the provider returns `permissions` as a Terraform list value, not a JSON string
