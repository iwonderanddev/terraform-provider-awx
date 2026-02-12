# Resource: awx_user

Manages AWX `users` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_user" "example" {
  username = "example"
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

- `username` (String, Required) Required. 150 characters or fewer. Letters, digits and @/./+/-/_ only.

### Optional

- `email` (String, Optional) Value for `email`.
- `first_name` (String, Optional) Value for `first_name`.
- `is_superuser` (Boolean, Optional, Computed) Designates that this user has all permissions without explicitly assigning them.
- `is_system_auditor` (Boolean, Optional, Computed) Controls whether `is_system_auditor` is enabled in AWX.
- `last_name` (String, Optional) Value for `last_name`.
- `password` (String, Optional, Sensitive, Write-Only) Field used to change the password.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_user.example 42
```

## Further Reading

- [AWX Users](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/users.html)
