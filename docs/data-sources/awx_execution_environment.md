# Data Source: awx_execution_environment

Reads AWX `execution_environments` objects.

## Example Usage

```hcl
data "awx_execution_environment" "example" {
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
- `description` (String, Read-Only) AWX value stored in `description`.
- `image` (String, Read-Only) The full image location, including the container registry, image name, and version tag.
- `name` (String, Read-Only) AWX value stored in `name`.
- `organization_id` (Number, Read-Only) The organization used to determine access to this execution environment.
- `pull` (String, Read-Only) Pull image before running?
  - `always` - Always pull container before running.
  - `missing` - Only pull the image if not present before running.
  - `never` - Never pull container before running.

## Further Reading

- [AWX Execution Environments](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/execution_environments.html)
