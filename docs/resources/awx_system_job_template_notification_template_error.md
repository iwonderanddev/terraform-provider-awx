# Resource: awx_system_job_template_notification_template_error

Manages `system_job_template_notification_template_error` relationships between `system_job_templates` and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_system_job_template_notification_template_error" "example" {
  parent_id = 12
  child_id  = 34
}
```

## Argument Reference

- `parent_id` (Number, Required) Parent object numeric ID.
- `child_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.

## Import

```bash
terraform import awx_system_job_template_notification_template_error.example 12:34
```
