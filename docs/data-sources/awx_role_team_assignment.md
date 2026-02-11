# Data Source: awx_role_team_assignment

Reads AWX `role_team_assignments` objects.

## Example Usage

```hcl
data "awx_role_team_assignment" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `object_ansible_id` (string)
- `object_id` (string)
- `role_definition_id` (integer)
- `team_id` (integer)
- `team_ansible_id` (string)
