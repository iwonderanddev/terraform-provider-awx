# Data Source: awx_organization

Reads AWX `organizations` objects.

## Example Usage

```hcl
data "awx_organization" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `default_environment_id` (Number, Read-Only) The default execution environment for jobs run by this organization.
- `description` (String, Read-Only) Value for `description`.
- `max_hosts` (Number, Read-Only) Maximum number of hosts allowed to be managed by this organization.
- `name` (String, Read-Only) Value for `name`.
- `opa_query_path` (String, Read-Only) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.

## Further Reading

- [AWX Organizations](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/organizations.html)
