# Data Source: awx_inventory_source

Reads AWX `inventory_sources` objects.

## Example Usage

```hcl
data "awx_inventory_source" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `credential_id` (integer)
- `description` (string)
- `enabled_value` (string)
- `enabled_var` (string)
- `execution_environment_id` (integer)
- `host_filter` (string)
- `inventory_id` (integer)
- `limit` (string)
- `name` (string)
- `overwrite` (boolean)
- `overwrite_vars` (boolean)
- `scm_branch` (string)
- `source` (string)
- `source_path` (string)
- `source_project_id` (integer)
- `source_vars` (string)
- `timeout` (integer)
- `update_cache_timeout` (integer)
- `update_on_launch` (boolean)
- `verbosity` (integer)
