# Data Source: awx_host

Reads AWX `hosts` objects.

## Example Usage

```hcl
data "awx_host" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `description` (string)
- `enabled` (boolean)
- `instance_id` (string)
- `inventory` (integer)
- `name` (string)
- `variables` (string)
