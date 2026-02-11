# Resource: awx_job_template_notification_template_error

Manages `job_template_notification_template_error` relationships between `job_templates` and `notification_templates` objects.

Breaking change: use `job_template_id` and `notification_template_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_job_template_notification_template_error" "example" {
  job_template_id = 12
  notification_template_id  = 34
}
```

## Argument Reference

- `job_template_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `job_template_id` (Number) Parent object numeric ID.
- `notification_template_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_job_template_notification_template_error.example 12:34
```
