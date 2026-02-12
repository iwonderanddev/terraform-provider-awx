# Resource: awx_inventory_label_association

Manages AWX associations between `inventories` and `labels` objects.

## Example Usage

```hcl
resource "awx_inventory_label_association" "example" {
  inventory_id = 12
  label_id = 34
}
```

## Schema

### Required

- `inventory_id` (Number, Required) Parent object numeric ID.
- `label_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `inventory_id` (Number, Read-Only) Parent object numeric ID.
- `label_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_inventory_label_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
- [AWX Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/job_templates.html)
