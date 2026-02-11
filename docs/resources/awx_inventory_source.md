# Resource: awx_inventory_source

Manages AWX `inventory_sources` objects.

## Example Usage

```hcl
resource "awx_inventory_source" "example" {
  inventory_id = awx_inventory.example.id
  name = "example"
  source = "example"
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `credential_id` (Optional) Cloud credential to use for inventory updates.
- `description` (Optional) Managed field from AWX OpenAPI schema.
- `enabled_value` (Optional) Only used when enabled_var is set. Value when the host is considered enabled. For example if enabled_var="status.power_state"and enabled_value="powered_on" with host variables:{   "status": {     "power_state": "powered_on",     "created": "2020-08-04T18:13:04+00:00",     "healthy": true    },    "name": "foobar",    "ip_address": "192.168.2.1"}The host would be marked enabled. If power_state where any value other than powered_on then the host would be disabled when imported. If the key is not found then the host will be enabled
- `enabled_var` (Optional) Retrieve the enabled state from the given dict of host variables. The enabled variable may be specified as "foo.bar", in which case the lookup will traverse into nested dicts, equivalent to: from_dict.get("foo", {}).get("bar", default)
- `execution_environment_id` (Optional) The container image to be used for execution.
- `host_filter` (Optional) This field is deprecated and will be removed in a future release. Regex where only matching hosts will be imported.
- `inventory_id` (Required) Managed field from AWX OpenAPI schema.
- `limit` (Optional) Enter host, group or pattern match
- `name` (Required) Managed field from AWX OpenAPI schema.
- `overwrite` (Optional, Computed) Overwrite local groups and hosts from remote inventory source.
- `overwrite_vars` (Optional, Computed) Overwrite local variables from remote inventory source.
- `scm_branch` (Optional) Inventory source SCM branch. Project default used if blank. Only allowed if project allow_override field is set to true.
- `source` (Required) * `azure_rm` - Microsoft Azure Resource Manager
  - `controller` - Red Hat Ansible Automation Platform
  - `ec2` - Amazon EC2
  - `gce` - Google Compute Engine
  - `insights` - Red Hat Insights
  - `openshift_virtualization` - OpenShift Virtualization
  - `openstack` - OpenStack
  - `rhv` - Red Hat Virtualization
  - `satellite6` - Red Hat Satellite 6
  - `terraform` - Terraform State
  - `vmware` - VMware vCenter
  - `scm` - Sourced from a Project
  - `constructed` - Template additional groups and hostvars at runtime
- `source_path` (Optional) Managed field from AWX OpenAPI schema.
- `source_project` (Optional) Project containing inventory file used as source.
- `source_vars` (Optional) Inventory source variables in YAML or JSON format.
- `timeout` (Optional, Computed) The amount of time (in seconds) to run before the task is canceled.
- `update_cache_timeout` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `update_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `verbosity` (Optional, Computed) * `0` - 0 (WARNING)
  - `1` - 1 (INFO)
  - `2` - 2 (DEBUG)

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_inventory_source.example 42
```
