# awx-native-object-field-values Delta

## ADDED Requirements

### Requirement: Survey spec relationship `spec` uses native structured values

The provider SHALL represent the `spec` argument on survey-spec relationship resources (`awx_job_template_survey_spec`, `awx_workflow_job_template_survey_spec`) as a Terraform structured object value consistent with other native complex fields. Configuration for `spec` MUST NOT accept JSON-encoded strings.

#### Scenario: Survey spec accepts object configuration

- **WHEN** a user sets `spec` on a survey-spec relationship resource using Terraform object/map syntax
- **THEN** Terraform accepts the value without requiring `jsonencode` and the provider sends the decoded object payload to AWX

#### Scenario: Survey spec rejects string JSON

- **WHEN** a user assigns a bare JSON string to `spec`
- **THEN** Terraform emits a type diagnostic because `spec` is not a string attribute
