# awx-single-object-resource-model Delta

## ADDED Requirements

### Requirement: Selected manifest array fields use native Terraform lists

The provider SHALL expose `awx_role_definition.permissions` as Terraform `list(string)` for resources and data sources when the AWX OpenAPI schema defines `permissions` as an array of strings. This field MUST NOT use JSON-string configuration or state encoding.

#### Scenario: Role definition permissions are a native list

- **WHEN** a managed resource schema includes `permissions` on `role_definitions`
- **THEN** the field is configured and stored as Terraform list data, not a JSON-encoded string
