# Resource: awx_organization_notification_template_error

Manages `organization_notification_template_error` relationships between `organizations`
and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_organization_notification_template_error" "example" {
  organization_id = 12
  notification_template_id  = 34
}
```

## Argument Reference

- `organization_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `organization_id` (Number) Parent object numeric ID.
- `notification_template_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_organization_notification_template_error.example \
  12:34
```
