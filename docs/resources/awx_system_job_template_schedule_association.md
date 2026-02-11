# Resource: awx_system_job_template_schedule_association

Manages `system_job_template_schedule_association` relationships between `system_job_templates` and `schedules` objects.

Breaking change: use `system_job_template_id` and `schedule_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_system_job_template_schedule_association" "example" {
  system_job_template_id = 12
  schedule_id  = 34
}
```

## Argument Reference

- `system_job_template_id` (Number, Required) Parent object numeric ID.
- `schedule_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `system_job_template_id` (Number) Parent object numeric ID.
- `schedule_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_system_job_template_schedule_association.example 12:34
```
