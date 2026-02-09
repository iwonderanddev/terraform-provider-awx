# Resource: awx_system_job_template_schedule_association

Manages `system_job_template_schedule_association` relationships between `system_job_templates` and `schedules` objects.

## Example Usage

```hcl
resource "awx_system_job_template_schedule_association" "example" {
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
terraform import awx_system_job_template_schedule_association.example 12:34
```
