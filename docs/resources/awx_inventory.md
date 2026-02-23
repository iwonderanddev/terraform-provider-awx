# Resource: awx_inventory

Manages AWX inventories used as host collections for job execution.

## Example Usage

### Standard inventory

```hcl
resource "awx_organization" "platform" {
  name = "platform"
}

resource "awx_inventory" "production" {
  name            = "production"
  organization_id = awx_organization.platform.id
  description     = "Primary production inventory"
}
```

### Smart inventory

```hcl
resource "awx_organization" "platform" {
  name = "platform"
}

resource "awx_inventory" "linux_hosts" {
  name            = "linux-hosts"
  organization_id = awx_organization.platform.id
  kind            = "smart"
  host_filter     = "ansible_facts__os_family=RedHat"
}
```

## Schema

### Qualifiers

- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Computed`: AWX sets the value during create or refresh.
- `Sensitive`: Terraform redacts the value in normal CLI output.
- `Write-Only`: Sent to AWX during create/update and not read back.

### Required

- `name` (String, Required) Inventory name shown in AWX.
- `organization_id` (Number, Required) Numeric ID of the organization that owns the inventory.

### Optional

- `description` (String, Optional) Optional inventory description for operators.
- `host_filter` (String, Optional) Host filter expression used by smart inventories.
- `kind` (String, Optional) Inventory behavior mode (empty string for normal, `smart`, or `constructed`).
- `opa_query_path` (String, Optional) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.
- `prevent_instance_group_fallback` (Boolean, Optional, Computed) Prevents AWX from automatically applying organization/global instance group fallback behavior.
- `variables` (String, Optional) JSON-encoded object of inventory-level variables.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_inventory.example 42
```

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
