# Data Source: awx_role_user_assignment

Reads AWX `role_user_assignments` objects.

## Example Usage

```hcl
data "awx_role_user_assignment" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `content_type` (string)
- `created` (string)
- `created_by` (integer)
- `id` (integer)
- `object_ansible_id` (string)
- `object_id` (string)
- `related` (object)
- `role_definition` (integer)
- `summary_fields` (object)
- `url` (string)
- `user` (integer)
- `user_ansible_id` (string)
