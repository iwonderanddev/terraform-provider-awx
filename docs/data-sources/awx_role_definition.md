# Data Source: awx_role_definition

Reads AWX `role_definitions` objects.

## Example Usage

```hcl
data "awx_role_definition" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `content_type` (string)
- `description` (string)
- `name` (string)
- `permissions` (array)
