# Resource: awx_project_notification_template_started

Manages `project_notification_template_started` relationships between `projects`
and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_project_notification_template_started" "example" {
  project_id = 12
  notification_template_id  = 34
}
```

## Argument Reference

- `project_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `project_id` (Number) Parent object numeric ID.
- `notification_template_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_project_notification_template_started.example \
  12:34
```
