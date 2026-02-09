# Data Source: awx_group

Reads AWX `groups` objects.

## Example Usage

```hcl
data "awx_group" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `description` (string)
- `inventory` (integer)
- `name` (string)
- `variables` (string)
