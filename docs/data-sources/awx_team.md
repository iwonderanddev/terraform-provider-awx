# Data Source: awx_team

Reads AWX `teams` objects.

## Example Usage

```hcl
data "awx_team" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `description` (string)
- `name` (string)
- `organization` (integer)
