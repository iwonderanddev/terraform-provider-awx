# Resource: awx_inventory_source_notification_template_started

Manages `inventory_source_notification_template_started` relationships between `inventory_sources` and `notification_templates` objects.

Breaking change: use `inventory_source_id` and `notification_template_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_inventory_source_notification_template_started" "example" {
  inventory_source_id = 12
  notification_template_id  = 34
}
```

## Argument Reference

- `inventory_source_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `inventory_source_id` (Number) Parent object numeric ID.
- `notification_template_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_inventory_source_notification_template_started.example 12:34
```
