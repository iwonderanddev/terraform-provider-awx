# Resource: awx_team_credential_association

Manages AWX associations between `teams` and `credentials` objects.

## Example Usage

```hcl
resource "awx_team_credential_association" "example" {
  team_id = 12
  credential_id = 34
}
```

## Schema

### Required

- `team_id` (Number, Required) Parent object numeric ID.
- `credential_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `team_id` (Number, Read-Only) Parent object numeric ID.
- `credential_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_team_credential_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Teams](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/teams.html)
- [AWX Credentials](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credentials.html)
