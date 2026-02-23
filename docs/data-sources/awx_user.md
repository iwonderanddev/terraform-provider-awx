# Data Source: awx_user

Reads AWX `users` objects.

## Example Usage

```hcl
data "awx_user" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `email` (String, Read-Only) Email address used for notifications and API identity context in AWX.
- `first_name` (String, Read-Only) Given name shown in the AWX UI for this user account.
- `is_superuser` (Boolean, Read-Only) Designates that this user has all permissions without explicitly assigning them.
- `is_system_auditor` (Boolean, Read-Only) Controls whether `is_system_auditor` is enabled in AWX.
- `last_name` (String, Read-Only) Family name shown in the AWX UI for this user account.
- `password` (String, Read-Only, Sensitive) Field used to change the password.
- `username` (String, Read-Only) Required. 150 characters or fewer. Letters, digits and @/./+/-/_ only.

## Further Reading

- [AWX Users](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/users.html)
