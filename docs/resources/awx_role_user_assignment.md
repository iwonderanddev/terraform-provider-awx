# Resource: awx_role_user_assignment

Manages AWX `role_user_assignments` objects.

This endpoint does not support in-place updates; Terraform replaces the resource when arguments change.

## Example Usage

### Basic configuration

```hcl
resource "awx_role_user_assignment" "example" {
  role_definition_id = awx_role_definition.example.id
  related = { key = "value" }
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

- `content_type` (String, Optional, Computed) String to use for references to this type from other models in the API.
- `created` (String, Optional, Computed) The date/time this resource was created.
- `created_by_id` (Number, Optional, Computed) The user who created this resource.
- `id` (Number, Optional, Computed) Numeric AWX value used for `id`.
- `object_ansible_id` (String, Optional) The resource id of the object this role applies to. An alternative to the object_id field.
- `object_id` (Number, Optional) The primary key of the object this assignment applies to; null value indicates system-wide assignment.
- `related` (Object, Optional, Computed) Object value for `related`.
- `summary_fields` (Object, Optional, Computed) Object value for `summary_fields`.
- `url` (String, Optional, Computed) AWX value stored in `url`.
- `user_id` (Number, Optional) Numeric ID of the related AWX user object.
- `user_ansible_id` (String, Optional) The resource ID of the user who will receive permissions from this assignment. An alternative to user field.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_role_user_assignment.example 42
```

## Further Reading

- [AWX Role-Based Access Controls](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/rbac.html)
