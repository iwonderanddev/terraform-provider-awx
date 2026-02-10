# awx-relationship-resources Specification

## Purpose
TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.
## Requirements
### Requirement: Explicit relationship resource modeling
The provider SHALL model AWX object associations with dedicated relationship resources when the association has independent lifecycle semantics.

#### Scenario: Manage user-to-team membership
- **WHEN** a user declares a team membership relationship resource
- **THEN** the provider creates or removes only that membership association in AWX

### Requirement: Relationship resource identity
Relationship resources SHALL expose stable endpoint-aligned identity:
- association resources use composite IDs based on related object identifiers.
- parent-scoped singleton relationship resources use parent-key IDs.

#### Scenario: Relationship state refresh
- **WHEN** the provider refreshes a relationship resource with ID `<left_id>:<right_id>`
- **THEN** the provider verifies the association exists and preserves the same composite ID format in state

#### Scenario: Singleton relationship state refresh
- **WHEN** the provider refreshes a parent-scoped singleton relationship resource with ID `<parent_id>`
- **THEN** the provider verifies the relationship exists and preserves the same parent-key ID format in state

