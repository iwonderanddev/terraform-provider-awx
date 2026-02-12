# Data Source: awx_host

Reads AWX `hosts` objects.

## Example Usage

```hcl
data "awx_host" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `description` (String, Read-Only) Value for `description`.
- `enabled` (Boolean, Read-Only) Is this host online and available for running jobs?
- `instance_id` (String, Read-Only) The value used by the remote inventory source to uniquely identify the host
- `inventory_id` (Number, Read-Only) Numeric ID of the related AWX inventory object.
- `name` (String, Read-Only) Value for `name`.
- `variables` (String, Read-Only) Host variables in JSON or YAML format.

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
