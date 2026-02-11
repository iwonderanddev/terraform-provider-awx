# Data Source: awx_inventory

Reads AWX `inventories` objects.

## Example Usage

```hcl
data "awx_inventory" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `description` (string)
- `host_filter` (string)
- `kind` (string)
- `name` (string)
- `opa_query_path` (string)
- `organization_id` (integer)
- `prevent_instance_group_fallback` (boolean)
- `variables` (string)
