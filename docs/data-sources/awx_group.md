# Data Source: awx_group

Reads AWX `groups` objects.

## Example Usage

```hcl
data "awx_group" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `description` (String, Read-Only) AWX value stored in `description`.
- `inventory_id` (Number, Read-Only) Numeric ID of the related AWX inventory object.
- `name` (String, Read-Only) AWX value stored in `name`.
- `variables` (String, Read-Only) Group variables in JSON or YAML format.

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
