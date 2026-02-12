# Resource: awx_job_template_notification_template_started

Manages AWX associations between `job_templates` and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_job_template_notification_template_started" "example" {
  job_template_id = 12
  notification_template_id = 34
}
```

## Schema

### Required

- `job_template_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `job_template_id` (Number, Read-Only) Parent object numeric ID.
- `notification_template_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_job_template_notification_template_started.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/job_templates.html)
- [AWX Notifications](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/notifications.html)
