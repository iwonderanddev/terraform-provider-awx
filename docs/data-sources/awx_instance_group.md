# Data Source: awx_instance_group

Reads AWX `instance_groups` objects.

## Example Usage

```hcl
data "awx_instance_group" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `credential_id` (Number, Read-Only) Numeric ID of the related AWX credential object.
- `is_container_group` (Boolean, Read-Only) Indicates whether instances in this group are containerized.Containerized groups have a designated Openshift or Kubernetes cluster.
- `max_concurrent_jobs` (Number, Read-Only) Maximum number of concurrent jobs to run on a group. When set to zero, no maximum is enforced.
- `max_forks` (Number, Read-Only) Maximum number of forks to execute concurrently on a group. When set to zero, no maximum is enforced.
- `name` (String, Read-Only) AWX value stored in `name`.
- `pod_spec_override` (String, Read-Only) AWX value stored in `pod_spec_override`.
- `policy_instance_list` (String, Read-Only) List of exact-match Instances that will be assigned to this group
- `policy_instance_minimum` (Number, Read-Only) Static minimum number of Instances that will be automatically assign to this group when new instances come online.
- `policy_instance_percentage` (Number, Read-Only) Minimum percentage of all instances that will be automatically assigned to this group when new instances come online.

## Further Reading

- [AWX Instance Groups](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/instance_groups.html)
