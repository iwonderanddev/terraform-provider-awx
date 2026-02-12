# Resource: awx_credential_input_source

Manages AWX `credential_input_sources` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_credential_input_source" "example" {
  input_field_name = "example"
  source_credential_id = 1
  target_credential_id = 1
}
```

## Schema

### Qualifiers

- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Computed`: AWX sets the value during create or refresh.
- `Sensitive`: Terraform redacts the value in normal CLI output.
- `Write-Only`: Sent to AWX during create/update and not read back.

### Required

- `input_field_name` (String, Required) Value for `input_field_name`.
- `source_credential_id` (Number, Required) Numeric ID of the related AWX source credential object.
- `target_credential_id` (Number, Required) Numeric ID of the related AWX target credential object.

### Optional

- `description` (String, Optional) Value for `description`.
- `metadata` (String, Optional, Computed) Value for `metadata`.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_credential_input_source.example 42
```

## Further Reading

- [AWX Secret Management System](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credential_plugins.html)
