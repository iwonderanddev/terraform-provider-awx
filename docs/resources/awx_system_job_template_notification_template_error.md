# Resource: awx_system_job_template_notification_template_error

Manages `system_job_template_notification_template_error` relationships between `system_job_templates`
and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_system_job_template_notification_template_error" "example" {
  system_job_template_id = 12
  notification_template_id  = 34
}
```

## Argument Reference

- `system_job_template_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `system_job_template_id` (Number) Parent object numeric ID.
- `notification_template_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_system_job_template_notification_template_error.example \
  12:34
```
