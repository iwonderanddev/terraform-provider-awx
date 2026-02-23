# Delta Specification: awx-relationship-resources

## MODIFIED Requirements

### Requirement: Explicit relationship resource modeling

The provider SHALL model AWX object associations with dedicated relationship
resources when the association has independent lifecycle semantics.
Relationship resources SHALL expose canonical object-specific identifier
arguments using explicit `_id` suffix names derived from relationship metadata
(for example, `team_id`, `user_id`) instead of generic directional names.

Relationship derivation SHALL support explicit child-collection alias mapping
when AWX exposes endpoint variants that represent the same managed child object.
When multiple endpoint paths resolve to the same relationship resource name,
the provider SHALL apply deterministic path-preference rules so the selected
path preserves association attach/detach semantics.

#### Scenario: Manage user-to-team membership

- **WHEN** a user declares a team membership relationship resource
- **THEN** the provider creates or removes only that membership association in
  AWX using canonical `team_id` and `user_id` arguments

#### Scenario: Configure survey spec relationship

- **WHEN** a user declares a survey spec relationship resource
- **THEN** the provider accepts canonical parent object identifier argument
  naming (for example, `job_template_id`) together with `spec`

#### Scenario: Legacy directional arguments are unsupported

- **WHEN** a user configures a relationship resource with non-canonical
  directional argument names
- **THEN** Terraform configuration validation fails because only canonical
  object-specific `*_id` arguments are supported

#### Scenario: Organization credential relationship uses galaxy endpoint

- **WHEN** relationship metadata is derived for organization-to-credential
  association paths
- **THEN** `organization_credential_association` resolves to
  `/api/v2/organizations/{id}/galaxy_credentials/` while still using canonical
  `organization_id` and `credential_id` arguments

#### Scenario: Deterministic path preference on name collision

- **WHEN** two eligible AWX relationship-like endpoints resolve to the same
  relationship resource name
- **THEN** derivation applies explicit path preference rules rather than
  first-seen lexical order

#### Scenario: Non-organization credential associations remain stable

- **WHEN** relationship metadata is derived for credential associations outside
  organizations (for example job templates, inventory sources, teams, users)
- **THEN** those relationships keep their existing endpoint mappings and
  canonical argument naming
