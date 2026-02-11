# Resource: awx_project_notification_template_error

Manages `project_notification_template_error` relationships between `projects` and `notification_templates` objects.

Breaking change: use `project_id` and `notification_template_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_project_notification_template_error" "example" {
  project_id = 12
  notification_template_id  = 34
}
```

## Argument Reference

- `project_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `project_id` (Number) Parent object numeric ID.
- `notification_template_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_project_notification_template_error.example 12:34
```
