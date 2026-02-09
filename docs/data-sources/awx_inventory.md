# Data Source: awx_inventory

Reads AWX `inventories` objects.

## Example Usage

```hcl
data "awx_inventory" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `description` (string)
- `host_filter` (string)
- `kind` (string)
- `name` (string)
- `opa_query_path` (string)
- `organization` (integer)
- `prevent_instance_group_fallback` (boolean)
- `variables` (string)
