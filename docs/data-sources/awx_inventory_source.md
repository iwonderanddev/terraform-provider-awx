# Data Source: awx_inventory_source

Reads AWX `inventory_sources` objects.

## Example Usage

```hcl
data "awx_inventory_source" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `credential_id` (Number, Read-Only) Numeric ID of the credential used to authenticate to the source system.
- `description` (String, Read-Only) Optional explanation of the source purpose.
- `enabled_value` (String, Read-Only) Only used when enabled_var is set. Value when the host is considered enabled. For example if enabled_var="status.power_state"and enabled_value="powered_on" with host variables:{   "status": {     "power_state": "powered_on",     "created": "2020-08-04T18:13:04+00:00",     "healthy": true    },    "name": "foobar",    "ip_address": "192.168.2.1"}The host would be marked enabled. If power_state where any value other than powered_on then the host would be disabled when imported. If the key is not found then the host will be enabled
- `enabled_var` (String, Read-Only) Retrieve the enabled state from the given dict of host variables. The enabled variable may be specified as "foo.bar", in which case the lookup will traverse into nested dicts, equivalent to: from_dict.get("foo", {}).get("bar", default)
- `execution_environment_id` (Number, Read-Only) The container image to be used for execution.
- `host_filter` (String, Read-Only) This field is deprecated and will be removed in a future release. Regex where only matching hosts will be imported.
- `inventory_id` (Number, Read-Only) Numeric ID of the destination inventory populated by this source.
- `limit` (String, Read-Only) Enter host, group or pattern match
- `name` (String, Read-Only) Inventory source name shown in AWX.
- `overwrite` (Boolean, Read-Only) When true, AWX replaces local groups/hosts using the synced source data.
- `overwrite_vars` (Boolean, Read-Only) When true, AWX replaces local variable values during sync.
- `scm_branch` (String, Read-Only) Inventory source SCM branch. Project default used if blank. Only allowed if project allow_override field is set to true.
- `source` (String, Read-Only) Dynamic inventory backend plugin (for example `ec2`, `gce`, `scm`, or `constructed`).
- `source_path` (String, Read-Only) Relative path in the source project that points to the inventory file.
- `source_project_id` (Number, Read-Only) Numeric ID of the AWX project used when `source = "scm"`.
- `source_vars` (String, Read-Only) JSON-encoded object or YAML string of source-specific plugin variables.
- `timeout` (Number, Read-Only) Maximum sync runtime in seconds before AWX cancels the update.
- `update_cache_timeout` (Number, Read-Only) Seconds AWX waits before allowing another automatic source refresh.
- `update_on_launch` (Boolean, Read-Only) When true, AWX refreshes this source before dependent launches.
- `verbosity` (Number, Read-Only) * `0` - 0 (WARNING)
  - `1` - 1 (INFO)
  - `2` - 2 (DEBUG)

## Further Reading

- [AWX Inventory Sources](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
- [AWX Inventory Sources](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html#inventory-sources)
