# Resource: awx_inventory

Manages AWX `inventories` objects.

## Example Usage

```hcl
resource "awx_inventory" "example" {
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
- `host_filter` (Optional) Filter that will be applied to the hosts of this inventory.
- `kind` (Optional) Kind of inventory being represented.
  - `` - Hosts have a direct link to this inventory.
  - `smart` - Hosts for inventory generated using the host_filter property.
  - `constructed` - Parse list of source inventories with the constructed inventory plugin.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `opa_query_path` (Optional) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.
- `organization` (Required) Organization containing this inventory.
- `prevent_instance_group_fallback` (Optional, Computed) If enabled, the inventory will prevent adding any organization instance groups to the list of preferred instances groups to run associated job templates on.If this setting is enabled and you provided an empty list, the global instance groups will be applied.
- `variables` (Optional) Inventory variables in JSON or YAML format.

## Attributes Reference

- `id` (String) Numeric AWX object identifier.

## Import

```bash
terraform import awx_inventory.example 42
```
