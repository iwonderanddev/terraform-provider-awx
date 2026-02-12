# Resource: awx_organization_user_association

Manages AWX associations between `organizations` and `users` objects.

## Example Usage

```hcl
resource "awx_organization_user_association" "example" {
  organization_id = 12
  user_id = 34
}
```

## Schema

### Required

- `organization_id` (Number, Required) Parent object numeric ID.
- `user_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `organization_id` (Number, Read-Only) Parent object numeric ID.
- `user_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_organization_user_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Organizations](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/organizations.html)
- [AWX Users](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/users.html)
