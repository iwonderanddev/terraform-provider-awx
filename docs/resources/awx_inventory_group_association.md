# Resource: awx_inventory_group_association

Manages `inventory_group_association` relationships between `inventories`
and `groups` objects.

## Example Usage

```hcl
resource "awx_inventory_group_association" "example" {
  inventory_id = 12
  group_id  = 34
}
```

## Argument Reference

- `inventory_id` (Number, Required) Parent object numeric ID.
- `group_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `inventory_id` (Number) Parent object numeric ID.
- `group_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_inventory_group_association.example \
  12:34
```
