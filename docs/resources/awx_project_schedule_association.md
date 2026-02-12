# Resource: awx_project_schedule_association

Manages AWX associations between `projects` and `schedules` objects.

## Example Usage

```hcl
resource "awx_project_schedule_association" "example" {
  project_id = 12
  schedule_id = 34
}
```

## Schema

### Required

- `project_id` (Number, Required) Parent object numeric ID.
- `schedule_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `project_id` (Number, Read-Only) Parent object numeric ID.
- `schedule_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_project_schedule_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Projects](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/projects.html)
- [AWX Schedules](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/scheduling.html)
