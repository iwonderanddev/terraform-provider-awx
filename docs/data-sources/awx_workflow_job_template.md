# Data Source: awx_workflow_job_template

Reads AWX `workflow_job_templates` objects.

## Example Usage

```hcl
data "awx_workflow_job_template" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `allow_simultaneous` (boolean)
- `ask_inventory_on_launch` (boolean)
- `ask_labels_on_launch` (boolean)
- `ask_limit_on_launch` (boolean)
- `ask_scm_branch_on_launch` (boolean)
- `ask_skip_tags_on_launch` (boolean)
- `ask_tags_on_launch` (boolean)
- `ask_variables_on_launch` (boolean)
- `description` (string)
- `extra_vars` (object)
- `inventory_id` (integer)
- `job_tags` (string)
- `limit` (string)
- `name` (string)
- `organization_id` (integer)
- `scm_branch` (string)
- `skip_tags` (string)
- `survey_enabled` (boolean)
- `webhook_credential` (integer)
- `webhook_service` (string)
