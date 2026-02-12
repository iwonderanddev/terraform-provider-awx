# Resource: awx_project_schedule_association

Manages `project_schedule_association` relationships between `projects`
and `schedules` objects.

## Example Usage

```hcl
resource "awx_project_schedule_association" "example" {
  project_id = 12
  schedule_id  = 34
}
```

## Argument Reference

- `project_id` (Number, Required) Parent object numeric ID.
- `schedule_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `project_id` (Number) Parent object numeric ID.
- `schedule_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_project_schedule_association.example \
  12:34
```
