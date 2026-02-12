# Resource: awx_workflow_job_template

Manages AWX workflow job templates that coordinate multi-step automation across linked nodes.

## AWX Concepts

Workflow job templates orchestrate multiple job templates with conditional paths. They are useful when releases or operations require ordered stages, approvals, and failure handling in one AWX launch artifact.

## Example Usage

### Workflow template for release orchestration

```hcl
resource "awx_workflow_job_template" "release" {
  name            = "release-workflow"
  organization_id = awx_organization.platform.id
  survey_enabled  = true
  ask_variables_on_launch = true
  extra_vars = {
    release_track = "stable"
  }
}
```

### Webhook-triggered workflow

```hcl
resource "awx_workflow_job_template" "from_gitlab" {
  name                  = "workflow-from-gitlab"
  organization_id       = awx_organization.platform.id
  webhook_service       = "gitlab"
  webhook_credential_id = awx_credential.gitlab_token.id
}
```

## Schema

### Qualifiers

- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Computed`: AWX sets the value during create or refresh.
- `Sensitive`: Terraform redacts the value in normal CLI output.
- `Write-Only`: Sent to AWX during create/update and not read back.

### Required

- `name` (String, Required) Human-readable workflow template name in AWX.

### Optional

- `allow_simultaneous` (Boolean, Optional, Computed) Controls whether `allow_simultaneous` is enabled in AWX.
- `ask_inventory_on_launch` (Boolean, Optional, Computed) When true, launch requests can override inventory.
- `ask_labels_on_launch` (Boolean, Optional, Computed) Controls whether `ask_labels_on_launch` is enabled in AWX.
- `ask_limit_on_launch` (Boolean, Optional, Computed) When true, launch requests can set a limit pattern.
- `ask_scm_branch_on_launch` (Boolean, Optional, Computed) Controls whether `ask_scm_branch_on_launch` is enabled in AWX.
- `ask_skip_tags_on_launch` (Boolean, Optional, Computed) Controls whether `ask_skip_tags_on_launch` is enabled in AWX.
- `ask_tags_on_launch` (Boolean, Optional, Computed) Controls whether `ask_tags_on_launch` is enabled in AWX.
- `ask_variables_on_launch` (Boolean, Optional, Computed) When true, launch requests can provide extra vars.
- `description` (String, Optional) Optional summary shown in AWX for operators.
- `extra_vars` (Object, Optional) Object of default extra vars used by the workflow launch.
- `inventory_id` (Number, Optional) Numeric ID of the inventory prompt default for workflow launches.
- `job_tags` (String, Optional) Value for `job_tags`.
- `limit` (String, Optional) Value for `limit`.
- `organization_id` (Number, Optional) Numeric ID of the owning organization.
- `scm_branch` (String, Optional) Value for `scm_branch`.
- `skip_tags` (String, Optional) Value for `skip_tags`.
- `survey_enabled` (Boolean, Optional, Computed) Enables survey prompts for workflow launches.
- `webhook_credential_id` (Number, Optional) Numeric ID of the credential used for webhook signature validation.
- `webhook_service` (String, Optional) Webhook provider accepted for workflow launches.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_workflow_job_template.example 42
```

## Further Reading

- [AWX Workflow Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/workflow_templates.html)
