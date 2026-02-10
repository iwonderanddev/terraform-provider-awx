# Data Source: awx_role

Reads AWX `roles` objects.

## Example Usage

```hcl
data "awx_role" "example" {
  id = "example"
}
```

## Argument Reference

- `id` (String, Optional) AWX object identifier used in the detail endpoint path.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (String) AWX detail-path identifier for this object.
- `description` (string)
- `id` (integer)
- `name` (string)
- `related` (string)
- `summary_fields` (string)
- `type` (string)
- `url` (string)
