# Resource: awx_inventory_source_notification_template_error

Manages `inventory_source_notification_template_error` relationships between `inventory_sources`
and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_inventory_source_notification_template_error" "example" {
  inventory_source_id = 12
  notification_template_id  = 34
}
```

## Argument Reference

- `inventory_source_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `inventory_source_id` (Number) Parent object numeric ID.
- `notification_template_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_inventory_source_notification_template_error.example \
  12:34
```
