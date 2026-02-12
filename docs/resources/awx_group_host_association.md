# Resource: awx_group_host_association

Manages AWX associations between `groups` and `hosts` objects.

## Example Usage

```hcl
resource "awx_group_host_association" "example" {
  group_id = 12
  host_id = 34
}
```

## Schema

### Required

- `group_id` (Number, Required) Parent object numeric ID.
- `host_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `group_id` (Number, Read-Only) Parent object numeric ID.
- `host_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_group_host_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
