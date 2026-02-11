# Resource: awx_role_team_assignment

Manages AWX `role_team_assignments` objects.

This endpoint does not support in-place updates; Terraform replaces the resource when arguments change.

## Example Usage

```hcl
resource "awx_role_team_assignment" "example" {
  role_definition_id = awx_role_definition.example.id
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `object_ansible_id` (Optional) The resource id of the object this role applies to. An alternative to the object_id field.
- `object_id` (Optional) The primary key of the object this assignment applies to; null value indicates system-wide assignment.
- `role_definition_id` (Required) The role definition which defines permissions conveyed by this assignment.
- `team_id` (Optional) Managed field from AWX OpenAPI schema.
- `team_ansible_id` (Optional) The resource ID of the team who will receive permissions from this assignment. An alternative to team field.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_role_team_assignment.example 42
```
