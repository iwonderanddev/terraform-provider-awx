# awx-sensitive-field-handling Specification

## Purpose
TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.
## Requirements
### Requirement: Sensitive schema declaration
The provider SHALL mark secret and password-like attributes as sensitive in Terraform schemas.

#### Scenario: Plan output for sensitive field
- **WHEN** a user sets a secret attribute in configuration
- **THEN** Terraform plan output redacts the sensitive value

### Requirement: Write-only secret state behavior
The provider SHALL avoid repopulating cleartext secret values into Terraform state during read operations. For write-only secret fields that are typed as Terraform object values, the provider SHALL preserve configured/planned state values as typed objects and SHALL NOT replace them from API read responses.

#### Scenario: Read after secret update
- **WHEN** a resource with secret inputs is refreshed after apply
- **THEN** the provider maintains state without exposing cleartext secret values returned or derived from API responses

#### Scenario: Read after write-only object secret update
- **WHEN** a write-only secret field with object semantics is configured and the resource is refreshed
- **THEN** state preserves the typed object secret value and does not repopulate from AWX read responses

