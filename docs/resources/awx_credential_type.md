# Resource: awx_credential_type

Manages AWX `credential_types` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_credential_type" "example" {
  kind = "example"
  name = "example"
  injectors = { key = "value" }
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

- `kind` (String, Required) Allowed values:
  - `cloud` - Cloud
  - `net` - Network
- `name` (String, Required) AWX value stored in `name`.

### Optional

- `description` (String, Optional) AWX value stored in `description`.
- `injectors` (Object, Optional, Computed) Terraform object defining credential type injectors. Refer to the documentation for expected structure.
- `inputs` (Object, Optional, Computed) Terraform object defining credential type inputs. Refer to the documentation for expected structure.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_credential_type.example 42
```

## Further Reading

- [AWX Credential Types](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credential_types.html)
