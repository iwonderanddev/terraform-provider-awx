# Data Source: awx_credential_type

Reads AWX `credential_types` objects.

## Example Usage

```hcl
data "awx_credential_type" "example" {
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
- `injectors` (Object, Read-Only) Terraform object defining credential type injectors. Refer to the documentation for expected structure.
- `inputs` (Object, Read-Only) Terraform object defining credential type inputs. Refer to the documentation for expected structure.
- `kind` (String, Read-Only) Allowed values:
  - `cloud` - Cloud
  - `net` - Network
- `name` (String, Read-Only) AWX value stored in `name`.

## Further Reading

- [AWX Credential Types](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credential_types.html)
