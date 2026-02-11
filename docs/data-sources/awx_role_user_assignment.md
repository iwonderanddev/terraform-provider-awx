# Data Source: awx_role_user_assignment

Reads AWX `role_user_assignments` objects.

## Example Usage

```hcl
data "awx_role_user_assignment" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `content_type` (string)
- `created` (string)
- `created_by` (integer)
- `id` (integer)
- `object_ansible_id` (string)
- `object_id` (string)
- `related` (object)
- `role_definition_id` (integer)
- `summary_fields` (object)
- `url` (string)
- `user_id` (integer)
- `user_ansible_id` (string)
