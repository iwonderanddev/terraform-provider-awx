# Resource: awx_constructed_inventory

Manages AWX `constructed_inventories` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_constructed_inventory" "example" {
  name = "example"
  organization_id = awx_organization.example.id
}
```

## Schema

### Qualifiers

- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Computed`: AWX sets the value during create or refresh.
- `Read-Only`: Cannot be set in configuration; Terraform records the value AWX returns.
- `Sensitive`: Terraform redacts the value in normal CLI output.
- `Write-Only`: Sent to AWX during create/update and not read back.

### Required

- `name` (String, Required) AWX value stored in `name`.
- `organization_id` (Number, Required) Organization containing this inventory.

### Optional

- `description` (String, Optional) AWX value stored in `description`.
- `limit` (String, Optional) The limit to restrict the returned hosts for the related auto-created inventory source, special to constructed inventory.
- `opa_query_path` (String, Optional) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.
- `prevent_instance_group_fallback` (Boolean, Optional, Computed) If enabled, the inventory will prevent adding any organization instance groups to the list of preferred instances groups to run associated job templates on.If this setting is enabled and you provided an empty list, the global instance groups will be applied.
- `source_vars` (String, Optional) The source_vars for the related auto-created inventory source, special to constructed inventory.
- `update_cache_timeout` (Number, Optional) The cache timeout for the related auto-created inventory source, special to constructed inventory
- `variables` (String, Optional) Inventory variables in JSON or YAML format.
- `verbosity` (Number, Optional) The verbosity level for the related auto-created inventory source, special to constructed inventory

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_constructed_inventory.example 42
```

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
