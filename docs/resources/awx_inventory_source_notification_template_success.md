# Resource: awx_inventory_source_notification_template_success

Manages AWX associations between `inventory_sources` and `notification_templates` objects.

## Example Usage

```hcl
resource "awx_inventory_source_notification_template_success" "example" {
  inventory_source_id = 12
  notification_template_id = 34
}
```

## Schema

### Required

- `inventory_source_id` (Number, Required) Parent object numeric ID.
- `notification_template_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `inventory_source_id` (Number, Read-Only) Parent object numeric ID.
- `notification_template_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_inventory_source_notification_template_success.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Inventory Sources](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
- [AWX Notifications](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/notifications.html)
