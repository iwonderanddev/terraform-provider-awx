# Data Source: awx_role_definition

Reads AWX `role_definitions` objects.

## Example Usage

```hcl
data "awx_role_definition" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `content_type` (String, Read-Only) String to use for references to this type from other models in the API.
- `description` (String, Read-Only) A description of this role.
- `name` (String, Read-Only) The name of this role.
- `permissions` (List of String, Read-Only) List of permission strings for `permissions`.

## Further Reading

- [AWX Role-Based Access Controls](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/rbac.html)
