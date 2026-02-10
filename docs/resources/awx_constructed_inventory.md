# Resource: awx_constructed_inventory

Manages AWX `constructed_inventories` objects.

## Example Usage

```hcl
resource "awx_constructed_inventory" "example" {
  name = "example"
  organization = 1
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `description` (Optional) Managed field from AWX OpenAPI schema.
- `limit` (Optional) The limit to restrict the returned hosts for the related auto-created inventory source, special to constructed inventory.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `opa_query_path` (Optional) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.
- `organization` (Required) Organization containing this inventory.
- `prevent_instance_group_fallback` (Optional, Computed) If enabled, the inventory will prevent adding any organization instance groups to the list of preferred instances groups to run associated job templates on.If this setting is enabled and you provided an empty list, the global instance groups will be applied.
- `source_vars` (Optional) The source_vars for the related auto-created inventory source, special to constructed inventory.
- `update_cache_timeout` (Optional) The cache timeout for the related auto-created inventory source, special to constructed inventory
- `variables` (Optional) Inventory variables in JSON or YAML format.
- `verbosity` (Optional) The verbosity level for the related auto-created inventory source, special to constructed inventory

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_constructed_inventory.example 42
```
