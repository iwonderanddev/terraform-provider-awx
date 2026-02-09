# Resource: awx_inventory_group_association

Manages `inventory_group_association` relationships between `inventories` and `groups` objects.

## Example Usage

```hcl
resource "awx_inventory_group_association" "example" {
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
terraform import awx_inventory_group_association.example 12:34
```
