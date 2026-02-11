# awx-native-object-field-values Specification

## Purpose
TBD - created by archiving change enforce-native-complex-fields. Update Purpose after archive.
## Requirements
### Requirement: Generated object fields are native Terraform object values
The provider SHALL represent every generated AWX field with manifest type `object` as a Terraform object value in both resource and data source schemas. Configuration for these fields MUST NOT accept JSON-encoded strings.

#### Scenario: Resource object field configuration is typed
- **WHEN** a user configures a generated object field on a resource
- **THEN** Terraform accepts object input and rejects string input with a type diagnostic

#### Scenario: Data source object field readback is typed
- **WHEN** a data source reads an AWX field typed as `object`
- **THEN** the provider returns that attribute as a Terraform object value rather than a JSON string

### Requirement: Job template extra_vars fields are object-typed
The provider SHALL expose `extra_vars` as a Terraform object field for both `awx_job_template` and `awx_workflow_job_template`, including resources and data sources.

#### Scenario: Resource extra_vars accepts object input
- **WHEN** a user sets `extra_vars` on a job template or workflow job template resource
- **THEN** the field is configured as a Terraform object value and converted by the provider for AWX transport

#### Scenario: Data source extra_vars normalizes structured string responses
- **WHEN** AWX returns `extra_vars` as JSON or YAML string content
- **THEN** the provider normalizes the value into a Terraform object and returns a diagnostic if the normalized root is not an object

