# Resource: awx_inventory_source_credential_association

Manages `inventory_source_credential_association` relationships between `inventory_sources` and `credentials` objects.

Breaking change: use `inventory_source_id` and `credential_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_inventory_source_credential_association" "example" {
  inventory_source_id = 12
  credential_id  = 34
}
```

## Argument Reference

- `inventory_source_id` (Number, Required) Parent object numeric ID.
- `credential_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `inventory_source_id` (Number) Parent object numeric ID.
- `credential_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_inventory_source_credential_association.example 12:34
```
