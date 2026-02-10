# Data Source: awx_instance_group

Reads AWX `instance_groups` objects.

## Example Usage

```hcl
data "awx_instance_group" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `credential` (integer)
- `is_container_group` (boolean)
- `max_concurrent_jobs` (integer)
- `max_forks` (integer)
- `name` (string)
- `pod_spec_override` (string)
- `policy_instance_list` (array)
- `policy_instance_minimum` (integer)
- `policy_instance_percentage` (integer)
