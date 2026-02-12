# Resource: awx_team_credential_association

Manages `team_credential_association` relationships between `teams`
and `credentials` objects.

## Example Usage

```hcl
resource "awx_team_credential_association" "example" {
  team_id = 12
  credential_id  = 34
}
```

## Argument Reference

- `team_id` (Number, Required) Parent object numeric ID.
- `credential_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `team_id` (Number) Parent object numeric ID.
- `credential_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_team_credential_association.example \
  12:34
```
