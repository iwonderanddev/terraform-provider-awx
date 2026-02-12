# Resource: awx_workflow_job_template_notification_template_error

Manages `workflow_job_template_notification_template_error` relationships between `workflow_job_templates`
and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_workflow_job_template_notification_template_error" "example" {
  workflow_job_template_id = 12
  notification_template_id  = 34
}
```

## Argument Reference

- `workflow_job_template_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `workflow_job_template_id` (Number) Parent object numeric ID.
- `notification_template_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_workflow_job_template_notification_template_error.example \
  12:34
```
