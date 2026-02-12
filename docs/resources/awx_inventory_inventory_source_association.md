# Resource: awx_inventory_inventory_source_association

Manages `inventory_inventory_source_association` relationships between `inventories`
and `inventory_sources` objects.

## Example Usage

```hcl
resource "awx_inventory_inventory_source_association" "example" {
  inventory_id = 12
  inventory_source_id  = 34
}
```

## Argument Reference

- `inventory_id` (Number, Required) Parent object numeric ID.
- `inventory_source_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `inventory_id` (Number) Parent object numeric ID.
- `inventory_source_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_inventory_inventory_source_association.example \
  12:34
```
