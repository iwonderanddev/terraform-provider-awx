# Data Source: awx_inventory

Reads AWX `inventories` objects.

## Example Usage

```hcl
data "awx_inventory" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `description` (String, Read-Only) Optional inventory description for operators.
- `host_filter` (String, Read-Only) Host filter expression used by smart inventories.
- `kind` (String, Read-Only) Inventory behavior mode (empty string for normal, `smart`, or `constructed`).
- `name` (String, Read-Only) Inventory name shown in AWX.
- `opa_query_path` (String, Read-Only) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.
- `organization_id` (Number, Read-Only) Numeric ID of the organization that owns the inventory.
- `prevent_instance_group_fallback` (Boolean, Read-Only) Prevents AWX from automatically applying organization/global instance group fallback behavior.
- `variables` (String, Read-Only) JSON-encoded object of inventory-level variables.

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
