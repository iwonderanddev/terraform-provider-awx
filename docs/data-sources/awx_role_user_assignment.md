# Data Source: awx_role_user_assignment

Reads AWX `role_user_assignments` objects.

## Example Usage

```hcl
data "awx_role_user_assignment" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `content_type` (String, Read-Only) String to use for references to this type from other models in the API.
- `created` (String, Read-Only) The date/time this resource was created.
- `created_by_id` (Number, Read-Only) The user who created this resource.
- `id` (Number, Read-Only) Numeric AWX value used for `id`.
- `object_ansible_id` (String, Read-Only) The resource id of the object this role applies to. An alternative to the object_id field.
- `object_id` (Number, Read-Only) The primary key of the object this assignment applies to; null value indicates system-wide assignment.
- `related` (Object, Read-Only) Object value for `related`.
- `role_definition_id` (Number, Read-Only) The role definition which defines permissions conveyed by this assignment.
- `summary_fields` (Object, Read-Only) Object value for `summary_fields`.
- `url` (String, Read-Only) AWX value stored in `url`.
- `user_id` (Number, Read-Only) Numeric ID of the related AWX user object.
- `user_ansible_id` (String, Read-Only) The resource ID of the user who will receive permissions from this assignment. An alternative to user field.

## Further Reading

- [AWX Role-Based Access Controls](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/rbac.html)
