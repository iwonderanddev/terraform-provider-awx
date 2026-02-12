# Resource: awx_host_group_association

Manages AWX associations between `hosts` and `groups` objects.

## Example Usage

```hcl
resource "awx_host_group_association" "example" {
  host_id = 12
  group_id = 34
}
```

## Schema

### Required

- `host_id` (Number, Required) Parent object numeric ID.
- `group_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `host_id` (Number, Read-Only) Parent object numeric ID.
- `group_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_host_group_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
