# awx-relationship-resources Specification

## Purpose

TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.
## Requirements
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

### Requirement: Relationship resource identity

Relationship resources SHALL expose stable endpoint-aligned identity:

- association resources use composite IDs based on related object identifiers.
- parent-scoped singleton relationship resources use parent-key IDs.

#### Scenario: Relationship state refresh

- **WHEN** the provider refreshes a relationship resource with ID `<left_id>:<right_id>`
- **THEN** the provider verifies the association exists and preserves the same composite ID format in state regardless of canonical argument naming

#### Scenario: Singleton relationship state refresh

- **WHEN** the provider refreshes a parent-scoped singleton relationship resource with ID `<resource_id>`
- **THEN** the provider verifies the relationship exists and preserves the same parent-key ID format in state

### Requirement: Workflow template node edge relationship modeling

The provider SHALL model AWX workflow template node edge endpoints as explicit
relationship resources when AWX exposes
`/api/v2/workflow_job_template_nodes/{id}/success_nodes/`,
`/api/v2/workflow_job_template_nodes/{id}/failure_nodes/`, or
`/api/v2/workflow_job_template_nodes/{id}/always_nodes/`.

These resources SHALL use the following Terraform resource names and canonical
arguments:

- `awx_workflow_job_template_node_success_node_association` with
  `workflow_job_template_node_id` and `success_node_id`
- `awx_workflow_job_template_node_failure_node_association` with
  `workflow_job_template_node_id` and `failure_node_id`
- `awx_workflow_job_template_node_always_node_association` with
  `workflow_job_template_node_id` and `always_node_id`

These resources SHALL preserve standard relationship lifecycle and identity
behavior, including composite import IDs in `<primary_id>:<related_id>` format.

#### Scenario: Configure workflow node success edge

- **WHEN** a user declares
  `awx_workflow_job_template_node_success_node_association`
- **THEN** the provider manages the AWX
  `/api/v2/workflow_job_template_nodes/{id}/success_nodes/` association using
  `workflow_job_template_node_id` and `success_node_id`

#### Scenario: Configure workflow node failure edge

- **WHEN** a user declares
  `awx_workflow_job_template_node_failure_node_association`
- **THEN** the provider manages the AWX
  `/api/v2/workflow_job_template_nodes/{id}/failure_nodes/` association using
  `workflow_job_template_node_id` and `failure_node_id`

#### Scenario: Configure workflow node always edge

- **WHEN** a user declares
  `awx_workflow_job_template_node_always_node_association`
- **THEN** the provider manages the AWX
  `/api/v2/workflow_job_template_nodes/{id}/always_nodes/` association using
  `workflow_job_template_node_id` and `always_node_id`

#### Scenario: Runtime workflow job node edges remain unmanaged

- **WHEN** relationship metadata is derived for
  `/api/v2/workflow_job_nodes/{id}/{success|failure|always}_nodes/`
- **THEN** the provider does not expose managed relationship resources for
  those runtime workflow job node endpoints

