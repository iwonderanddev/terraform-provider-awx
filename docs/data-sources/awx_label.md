# Data Source: awx_label

Reads AWX `labels` objects.

## Example Usage

```hcl
data "awx_label" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `name` (string)
- `organization_id` (integer)
