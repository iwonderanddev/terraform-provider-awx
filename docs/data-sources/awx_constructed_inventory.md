# Data Source: awx_constructed_inventory

Reads AWX `constructed_inventories` objects.

## Example Usage

```hcl
data "awx_constructed_inventory" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `description` (String, Read-Only) AWX value stored in `description`.
- `limit` (String, Read-Only) The limit to restrict the returned hosts for the related auto-created inventory source, special to constructed inventory.
- `name` (String, Read-Only) AWX value stored in `name`.
- `opa_query_path` (String, Read-Only) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.
- `organization_id` (Number, Read-Only) Organization containing this inventory.
- `prevent_instance_group_fallback` (Boolean, Read-Only) If enabled, the inventory will prevent adding any organization instance groups to the list of preferred instances groups to run associated job templates on.If this setting is enabled and you provided an empty list, the global instance groups will be applied.
- `source_vars` (String, Read-Only) The source_vars for the related auto-created inventory source, special to constructed inventory.
- `update_cache_timeout` (Number, Read-Only) The cache timeout for the related auto-created inventory source, special to constructed inventory
- `variables` (String, Read-Only) Inventory variables in JSON or YAML format.
- `verbosity` (Number, Read-Only) The verbosity level for the related auto-created inventory source, special to constructed inventory

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
