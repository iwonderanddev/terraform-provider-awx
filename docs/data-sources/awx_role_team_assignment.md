# Data Source: awx_role_team_assignment

Reads AWX `role_team_assignments` objects.

## Example Usage

```hcl
data "awx_role_team_assignment" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `object_ansible_id` (string)
- `object_id` (string)
- `role_definition` (integer)
- `team` (integer)
- `team_ansible_id` (string)
