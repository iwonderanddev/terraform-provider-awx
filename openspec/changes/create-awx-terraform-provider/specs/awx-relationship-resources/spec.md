## ADDED Requirements

### Requirement: Explicit relationship resource modeling
The provider SHALL model AWX object associations with dedicated relationship resources when the association has independent lifecycle semantics.

#### Scenario: Manage user-to-team membership
- **WHEN** a user declares a team membership relationship resource
- **THEN** the provider creates or removes only that membership association in AWX

### Requirement: Relationship resource identity
Relationship resources SHALL expose stable identity using composite IDs based on related object identifiers.

#### Scenario: Relationship state refresh
- **WHEN** the provider refreshes a relationship resource with ID `<left_id>:<right_id>`
- **THEN** the provider verifies the association exists and preserves the same composite ID format in state
