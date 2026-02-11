# Resource: awx_inventory_label_association

Manages `inventory_label_association` relationships between `inventories` and `labels` objects.

Breaking change: use `inventory_id` and `label_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_inventory_label_association" "example" {
  inventory_id = 12
  label_id  = 34
}
```

## Argument Reference

- `inventory_id` (Number, Required) Parent object numeric ID.
- `label_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `inventory_id` (Number) Parent object numeric ID.
- `label_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_inventory_label_association.example 12:34
```
