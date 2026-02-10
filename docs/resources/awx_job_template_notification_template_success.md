# Resource: awx_job_template_notification_template_success

Manages `job_template_notification_template_success` relationships between `job_templates` and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_job_template_notification_template_success" "example" {
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
terraform import awx_job_template_notification_template_success.example 12:34
```
