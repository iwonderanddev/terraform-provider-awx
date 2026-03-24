# Resource: awx_instance_group

Manages AWX `instance_groups` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_instance_group" "example" {
  name = "example"
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

### Optional

- `credential_id` (Number, Optional) Numeric ID of the related AWX credential object.
- `is_container_group` (Boolean, Optional) Indicates whether instances in this group are containerized.Containerized groups have a designated Openshift or Kubernetes cluster.
- `max_concurrent_jobs` (Number, Optional, Computed) Maximum number of concurrent jobs to run on a group. When set to zero, no maximum is enforced.
- `max_forks` (Number, Optional, Computed) Maximum number of forks to execute concurrently on a group. When set to zero, no maximum is enforced.
- `pod_spec_override` (String, Optional) AWX value stored in `pod_spec_override`.
- `policy_instance_list` (String, Optional) List of exact-match Instances that will be assigned to this group
- `policy_instance_minimum` (Number, Optional, Computed) Static minimum number of Instances that will be automatically assign to this group when new instances come online.
- `policy_instance_percentage` (Number, Optional, Computed) Minimum percentage of all instances that will be automatically assigned to this group when new instances come online.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_instance_group.example 42
```

## Further Reading

- [AWX Instance Groups](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/instance_groups.html)
