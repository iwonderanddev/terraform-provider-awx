# Data Source: awx_team

Reads AWX `teams` objects.

## Example Usage

```hcl
data "awx_team" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `description` (string)
- `name` (string)
- `organization` (integer)
