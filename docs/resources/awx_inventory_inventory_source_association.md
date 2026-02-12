# Resource: awx_inventory_inventory_source_association

Manages AWX associations between `inventories` and `inventory_sources` objects.

## Example Usage

```hcl
resource "awx_inventory_inventory_source_association" "example" {
  inventory_id = 12
  inventory_source_id = 34
}
```

## Schema

### Required

- `inventory_id` (Number, Required) Parent object numeric ID.
- `inventory_source_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `inventory_id` (Number, Read-Only) Parent object numeric ID.
- `inventory_source_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_inventory_inventory_source_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
