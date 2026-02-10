# Resource: awx_organization_notification_template_approval

Manages `organization_notification_template_approval` relationships between `organizations` and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_organization_notification_template_approval" "example" {
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
terraform import awx_organization_notification_template_approval.example 12:34
```
