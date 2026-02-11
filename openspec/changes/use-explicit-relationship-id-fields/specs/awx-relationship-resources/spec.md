## MODIFIED Requirements

### Requirement: Explicit relationship resource modeling
The provider SHALL model AWX object associations with dedicated relationship resources when the association has independent lifecycle semantics. Relationship resources SHALL expose canonical object-specific identifier arguments using explicit `_id` suffix names derived from relationship metadata (for example, `team_id`, `user_id`) instead of generic directional names.

#### Scenario: Manage user-to-team membership
- **WHEN** a user declares a team membership relationship resource
- **THEN** the provider creates or removes only that membership association in AWX using canonical `team_id` and `user_id` arguments

#### Scenario: Configure survey spec relationship
- **WHEN** a user declares a survey spec relationship resource
- **THEN** the provider accepts canonical parent object identifier argument naming (for example, `job_template_id`) together with `spec`

#### Scenario: Legacy directional arguments are unsupported
- **WHEN** a user configures a relationship resource with `parent_id` or `child_id`
- **THEN** Terraform configuration validation fails because only canonical object-specific `*_id` arguments are supported

### Requirement: Relationship resource identity
Relationship resources SHALL expose stable endpoint-aligned identity:
- association resources use composite IDs based on related object identifiers.
- parent-scoped singleton relationship resources use parent-key IDs.

#### Scenario: Relationship state refresh
- **WHEN** the provider refreshes a relationship resource with ID `<left_id>:<right_id>`
- **THEN** the provider verifies the association exists and preserves the same composite ID format in state regardless of canonical argument naming

#### Scenario: Singleton relationship state refresh
- **WHEN** the provider refreshes a parent-scoped singleton relationship resource with ID `<parent_id>`
- **THEN** the provider verifies the relationship exists and preserves the same parent-key ID format in state
