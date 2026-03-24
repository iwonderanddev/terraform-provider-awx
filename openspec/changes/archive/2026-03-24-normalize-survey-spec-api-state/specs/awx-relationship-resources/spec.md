# awx-relationship-resources Delta

## MODIFIED Requirements

### Requirement: Explicit relationship resource modeling

The provider SHALL model AWX object associations with dedicated relationship
resources when the association has independent lifecycle semantics.
Relationship resources SHALL expose canonical object-specific identifier
arguments using explicit `_id` suffix names derived from relationship metadata
(for example, `team_id`, `user_id`) instead of generic directional names.

For survey-spec relationship resources (`awx_job_template_survey_spec`,
`awx_workflow_job_template_survey_spec`):

- On **create**, after a successful API create of the survey, the provider SHALL
  write the **configured** `spec` value to Terraform state (the same value sent
  in the POST body), so post-apply state matches the plan’s Dynamic/tuple typing.
- On **read** and **refresh**, the provider SHALL normalize AWX `survey_spec`
  GET bodies before converting the root object to Terraform `spec` state: each
  element of the `spec` question array SHALL be represented as an object whose
  attributes are the **union of keys present on at least one question in that
  GET response**, with absent keys on a given question set to null. The provider
  SHALL NOT inject attributes onto every question that never appear in the API
  or configuration for that survey.

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

#### Scenario: Survey spec create aligns state with plan

- **WHEN** a user applies a new survey-spec relationship resource
- **THEN** after a successful create, Terraform state for `spec` matches the
  configured value used for the POST body
- **AND** `terraform apply` does not fail with a post-apply inconsistency for
  `spec` solely because the AWX GET response encodes the survey differently

#### Scenario: Survey spec read normalizes sparse question objects

- **WHEN** the provider reads survey specification from AWX on refresh
- **AND** some questions omit keys that other questions in the same array include
- **THEN** the provider normalizes each question object to the union of keys
  observed in that GET payload before writing `spec` to state

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
- **THEN** derivation applies explicit path-preference rules rather than
  first-seen lexical order

#### Scenario: Non-organization credential associations remain stable

- **WHEN** relationship metadata is derived for credential associations outside
  organizations (for example job templates, inventory sources, teams, users)
- **THEN** those relationships keep their existing endpoint mappings and
  canonical argument naming
