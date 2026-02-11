# Resource: awx_instance_group

Manages AWX `instance_groups` objects.

## Example Usage

```hcl
resource "awx_instance_group" "example" {
  name = "example"
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `credential_id` (Optional) Managed field from AWX OpenAPI schema.
- `is_container_group` (Optional) Indicates whether instances in this group are containerized.Containerized groups have a designated Openshift or Kubernetes cluster.
- `max_concurrent_jobs` (Optional, Computed) Maximum number of concurrent jobs to run on a group. When set to zero, no maximum is enforced.
- `max_forks` (Optional, Computed) Maximum number of forks to execute concurrently on a group. When set to zero, no maximum is enforced.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `pod_spec_override` (Optional) Managed field from AWX OpenAPI schema.
- `policy_instance_list` (Optional) List of exact-match Instances that will be assigned to this group
- `policy_instance_minimum` (Optional, Computed) Static minimum number of Instances that will be automatically assign to this group when new instances come online.
- `policy_instance_percentage` (Optional, Computed) Minimum percentage of all instances that will be automatically assigned to this group when new instances come online.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_instance_group.example 42
```
