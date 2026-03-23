# awx-provider-documentation-and-examples Delta

## ADDED Requirements

### Requirement: Native structured examples for survey spec and role definition permissions

Generated documentation for `awx_job_template_survey_spec`, `awx_workflow_job_template_survey_spec`, and `awx_role_definition` SHALL show `spec` and `permissions` using Terraform structured syntax. Examples MUST NOT rely on `jsonencode` for these attributes.

#### Scenario: Survey spec documentation uses object syntax

- **WHEN** documentation is generated for a survey-spec relationship resource
- **THEN** `## Example Usage` shows `spec` as a Terraform object or map literal without `jsonencode`

#### Scenario: Role definition documentation uses list syntax

- **WHEN** documentation is generated for `awx_role_definition`
- **THEN** examples show `permissions` as a `list(string)` literal without `jsonencode`
