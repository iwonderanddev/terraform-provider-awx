# Resource: awx_job_template_notification_template_success

Manages `job_template_notification_template_success` relationships between `job_templates`
and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_job_template_notification_template_success" "example" {
  job_template_id = 12
  notification_template_id  = 34
}
```

## Argument Reference

- `job_template_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `job_template_id` (Number) Parent object numeric ID.
- `notification_template_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_job_template_notification_template_success.example \
  12:34
```
