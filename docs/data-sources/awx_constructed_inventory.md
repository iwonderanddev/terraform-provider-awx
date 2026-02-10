# Data Source: awx_constructed_inventory

Reads AWX `constructed_inventories` objects.

## Example Usage

```hcl
data "awx_constructed_inventory" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `description` (string)
- `limit` (string)
- `name` (string)
- `opa_query_path` (string)
- `organization` (integer)
- `prevent_instance_group_fallback` (boolean)
- `source_vars` (string)
- `update_cache_timeout` (integer)
- `variables` (string)
- `verbosity` (integer)
