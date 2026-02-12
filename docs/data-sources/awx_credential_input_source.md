# Data Source: awx_credential_input_source

Reads AWX `credential_input_sources` objects.

## Example Usage

```hcl
data "awx_credential_input_source" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `description` (String, Read-Only) Value for `description`.
- `input_field_name` (String, Read-Only) Value for `input_field_name`.
- `metadata` (String, Read-Only) Value for `metadata`.
- `source_credential_id` (Number, Read-Only) Numeric ID of the related AWX source credential object.
- `target_credential_id` (Number, Read-Only) Numeric ID of the related AWX target credential object.

## Further Reading

- [AWX Secret Management System](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credential_plugins.html)
