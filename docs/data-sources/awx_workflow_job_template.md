# Data Source: awx_workflow_job_template

Reads AWX `workflow_job_templates` objects.

## Example Usage

```hcl
data "awx_workflow_job_template" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `allow_simultaneous` (Boolean, Read-Only) Controls whether `allow_simultaneous` is enabled in AWX.
- `ask_inventory_on_launch` (Boolean, Read-Only) When true, launch requests can override inventory.
- `ask_labels_on_launch` (Boolean, Read-Only) Controls whether `ask_labels_on_launch` is enabled in AWX.
- `ask_limit_on_launch` (Boolean, Read-Only) When true, launch requests can set a limit pattern.
- `ask_scm_branch_on_launch` (Boolean, Read-Only) Controls whether `ask_scm_branch_on_launch` is enabled in AWX.
- `ask_skip_tags_on_launch` (Boolean, Read-Only) Controls whether `ask_skip_tags_on_launch` is enabled in AWX.
- `ask_tags_on_launch` (Boolean, Read-Only) Controls whether `ask_tags_on_launch` is enabled in AWX.
- `ask_variables_on_launch` (Boolean, Read-Only) When true, launch requests can provide extra vars.
- `description` (String, Read-Only) Optional summary shown in AWX for operators.
- `extra_vars` (Object, Read-Only) Object of default extra vars used by the workflow launch.
- `inventory_id` (Number, Read-Only) Numeric ID of the inventory prompt default for workflow launches.
- `job_tags` (String, Read-Only) Value for `job_tags`.
- `limit` (String, Read-Only) Value for `limit`.
- `name` (String, Read-Only) Human-readable workflow template name in AWX.
- `organization_id` (Number, Read-Only) Numeric ID of the owning organization.
- `scm_branch` (String, Read-Only) Value for `scm_branch`.
- `skip_tags` (String, Read-Only) Value for `skip_tags`.
- `survey_enabled` (Boolean, Read-Only) Enables survey prompts for workflow launches.
- `webhook_credential_id` (Number, Read-Only) Numeric ID of the credential used for webhook signature validation.
- `webhook_service` (String, Read-Only) Webhook provider accepted for workflow launches.

## Further Reading

- [AWX Workflow Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/workflow_templates.html)
