# Data Source: awx_role

Reads AWX `roles` objects.

## Example Usage

```hcl
data "awx_role" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `description` (string)
- `id` (integer)
- `name` (string)
- `related` (string)
- `summary_fields` (string)
- `type` (string)
- `url` (string)
