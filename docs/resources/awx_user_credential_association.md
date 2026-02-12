# Resource: awx_user_credential_association

Manages AWX associations between `users` and `credentials` objects.

## Example Usage

```hcl
resource "awx_user_credential_association" "example" {
  user_id = 12
  credential_id = 34
}
```

## Schema

### Required

- `user_id` (Number, Required) Parent object numeric ID.
- `credential_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `user_id` (Number, Read-Only) Parent object numeric ID.
- `credential_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_user_credential_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Users](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/users.html)
- [AWX Credentials](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credentials.html)
