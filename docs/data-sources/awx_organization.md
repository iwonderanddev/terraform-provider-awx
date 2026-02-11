# Data Source: awx_organization

Reads AWX `organizations` objects.

## Example Usage

```hcl
data "awx_organization" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `default_environment_id` (integer)
- `description` (string)
- `max_hosts` (integer)
- `name` (string)
- `opa_query_path` (string)
