# Resource: awx_inventory_source

Manages dynamic inventory sources that synchronize host data into AWX inventories.

## AWX Concepts

Inventory sources connect AWX inventories to external systems such as cloud APIs or SCM files. AWX can refresh these sources on launch and control overwrite behavior for imported groups, hosts, and variables.

## Example Usage

### EC2 dynamic inventory source

This example assumes an existing AWS credential named `aws-inventory`.

```hcl
resource "awx_organization" "platform" {
  name = "platform"
}

resource "awx_inventory" "production" {
  name            = "production"
  organization_id = awx_organization.platform.id
}

data "awx_credential" "aws_inventory" {
  name = "aws-inventory"
}

resource "awx_inventory_source" "aws_dynamic" {
  name             = "aws-dynamic"
  inventory_id     = awx_inventory.production.id
  source           = "ec2"
  credential_id    = data.awx_credential.aws_inventory.id
  update_on_launch = true
  overwrite        = true
}
```

### SCM-backed inventory source

```hcl
resource "awx_organization" "platform" {
  name = "platform"
}

resource "awx_inventory" "production" {
  name            = "production"
  organization_id = awx_organization.platform.id
}

resource "awx_project" "inventory_definitions" {
  name            = "inventory-definitions"
  organization_id = awx_organization.platform.id
  scm_type        = "git"
  scm_url         = "https://github.com/example/inventory-definitions.git"
}

resource "awx_inventory_source" "scm_hosts" {
  name              = "scm-hosts"
  inventory_id      = awx_inventory.production.id
  source            = "scm"
  source_project_id = awx_project.inventory_definitions.id
  source_path       = "inventories/prod.yml"
  update_on_launch  = true
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

- `inventory_id` (Number, Required) Numeric ID of the destination inventory populated by this source.
- `name` (String, Required) Inventory source name shown in AWX.
- `source` (String, Required) Dynamic inventory backend plugin (for example `ec2`, `gce`, `scm`, or `constructed`).

### Optional

- `credential_id` (Number, Optional) Numeric ID of the credential used to authenticate to the source system.
- `description` (String, Optional) Optional explanation of the source purpose.
- `enabled_value` (String, Optional) Only used when enabled_var is set. Value when the host is considered enabled. For example if enabled_var="status.power_state"and enabled_value="powered_on" with host variables:{   "status": {     "power_state": "powered_on",     "created": "2020-08-04T18:13:04+00:00",     "healthy": true    },    "name": "foobar",    "ip_address": "192.168.2.1"}The host would be marked enabled. If power_state where any value other than powered_on then the host would be disabled when imported. If the key is not found then the host will be enabled
- `enabled_var` (String, Optional) Retrieve the enabled state from the given dict of host variables. The enabled variable may be specified as "foo.bar", in which case the lookup will traverse into nested dicts, equivalent to: from_dict.get("foo", {}).get("bar", default)
- `execution_environment_id` (Number, Optional) The container image to be used for execution.
- `host_filter` (String, Optional) This field is deprecated and will be removed in a future release. Regex where only matching hosts will be imported.
- `limit` (String, Optional) Enter host, group or pattern match
- `overwrite` (Boolean, Optional, Computed) When true, AWX replaces local groups/hosts using the synced source data.
- `overwrite_vars` (Boolean, Optional, Computed) When true, AWX replaces local variable values during sync.
- `scm_branch` (String, Optional) Inventory source SCM branch. Project default used if blank. Only allowed if project allow_override field is set to true.
- `source_path` (String, Optional) Relative path in the source project that points to the inventory file.
- `source_project_id` (Number, Optional) Numeric ID of the AWX project used when `source = "scm"`.
- `source_vars` (String, Optional) JSON-encoded object or YAML string of source-specific plugin variables.
- `timeout` (Number, Optional, Computed) Maximum sync runtime in seconds before AWX cancels the update.
- `update_cache_timeout` (Number, Optional, Computed) Seconds AWX waits before allowing another automatic source refresh.
- `update_on_launch` (Boolean, Optional, Computed) When true, AWX refreshes this source before dependent launches.
- `verbosity` (Number, Optional, Computed) Allowed values:
  - `0` - 0 (WARNING)
  - `1` - 1 (INFO)
  - `2` - 2 (DEBUG)

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_inventory_source.example 42
```

## Further Reading

- [AWX Inventory Sources](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
- [AWX Inventory Sources](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html#inventory-sources)
