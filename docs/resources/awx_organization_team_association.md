# Resource: awx_organization_team_association

Manages AWX associations between `organizations` and `teams` objects.

## Example Usage

```hcl
resource "awx_organization_team_association" "example" {
  organization_id = 12
  team_id = 34
}
```

## Schema

### Required

- `organization_id` (Number, Required) Parent object numeric ID.
- `team_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `organization_id` (Number, Read-Only) Parent object numeric ID.
- `team_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_organization_team_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Organizations](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/organizations.html)
- [AWX Teams](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/teams.html)
