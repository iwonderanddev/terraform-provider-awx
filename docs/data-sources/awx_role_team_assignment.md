# Data Source: awx_role_team_assignment

Reads AWX `role_team_assignments` objects.

## Example Usage

```hcl
data "awx_role_team_assignment" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `object_ansible_id` (String, Read-Only) The resource id of the object this role applies to. An alternative to the object_id field.
- `object_id` (String, Read-Only) The primary key of the object this assignment applies to; null value indicates system-wide assignment.
- `role_definition_id` (Number, Read-Only) The role definition which defines permissions conveyed by this assignment.
- `team_id` (Number, Read-Only) Numeric ID of the related AWX team object.
- `team_ansible_id` (String, Read-Only) The resource ID of the team who will receive permissions from this assignment. An alternative to team field.

## Further Reading

- [AWX Role-Based Access Controls](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/rbac.html)
