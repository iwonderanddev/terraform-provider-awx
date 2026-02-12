# Resource: awx_system_job_template_schedule_association

Manages AWX associations between `system_job_templates` and `schedules` objects.

## Example Usage

```hcl
resource "awx_system_job_template_schedule_association" "example" {
  system_job_template_id = 12
  schedule_id = 34
}
```

## Schema

### Required

- `system_job_template_id` (Number, Required) Parent object numeric ID.
- `schedule_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `system_job_template_id` (Number, Read-Only) Parent object numeric ID.
- `schedule_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_system_job_template_schedule_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Jobs](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/jobs.html)
- [AWX Schedules](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/scheduling.html)
