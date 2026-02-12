# Resource: awx_team_user_association

Manages AWX associations between `teams` and `users` objects.

## Example Usage

```hcl
resource "awx_team_user_association" "example" {
  team_id = 12
  user_id = 34
}
```

## Schema

### Required

- `team_id` (Number, Required) Parent object numeric ID.
- `user_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `team_id` (Number, Read-Only) Parent object numeric ID.
- `user_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_team_user_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Teams](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/teams.html)
- [AWX Users](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/users.html)
