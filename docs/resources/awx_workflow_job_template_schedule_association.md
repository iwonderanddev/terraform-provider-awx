# Resource: awx_workflow_job_template_schedule_association

Manages `workflow_job_template_schedule_association` relationships between `workflow_job_templates` and `schedules` objects.

Breaking change: use `workflow_job_template_id` and `schedule_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_workflow_job_template_schedule_association" "example" {
  workflow_job_template_id = 12
  schedule_id  = 34
}
```

## Argument Reference

- `workflow_job_template_id` (Number, Required) Parent object numeric ID.
- `schedule_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `workflow_job_template_id` (Number) Parent object numeric ID.
- `schedule_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_workflow_job_template_schedule_association.example 12:34
```
