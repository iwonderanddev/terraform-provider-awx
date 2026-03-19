# Delta Specification: awx-relationship-resources

## ADDED Requirements

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
