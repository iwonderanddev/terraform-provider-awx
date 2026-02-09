# Data Source: awx_user

Reads AWX `users` objects.

## Example Usage

```hcl
data "awx_user" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `email` (string)
- `first_name` (string)
- `is_superuser` (boolean)
- `is_system_auditor` (boolean)
- `last_name` (string)
- `password` (string, Sensitive)
- `username` (string)
