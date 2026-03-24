# Resource: awx_role_team_assignment

Manages AWX `role_team_assignments` objects.

This endpoint does not support in-place updates; Terraform replaces the resource when arguments change.

## Example Usage

### Basic configuration

```hcl
resource "awx_role_team_assignment" "example" {
  role_definition_id = awx_role_definition.example.id
}
```

## Schema

### Qualifiers

- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Computed`: AWX sets the value during create or refresh.
- `Read-Only`: Cannot be set in configuration; Terraform records the value AWX returns.
- `Sensitive`: Terraform redacts the value in normal CLI output.
- `Write-Only`: Sent to AWX during create/update and not read back.

### Required

- `role_definition_id` (Number, Required) The role definition which defines permissions conveyed by this assignment.

### Optional

- `object_ansible_id` (String, Optional) The resource id of the object this role applies to. An alternative to the object_id field.
- `object_id` (Number, Optional) The primary key of the object this assignment applies to; null value indicates system-wide assignment.
- `team_id` (Number, Optional) Numeric ID of the related AWX team object.
- `team_ansible_id` (String, Optional) The resource ID of the team who will receive permissions from this assignment. An alternative to team field.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_role_team_assignment.example 42
```

## Further Reading

- [AWX Role-Based Access Controls](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/rbac.html)
