# Resource: awx_execution_environment

Manages AWX `execution_environments` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_execution_environment" "example" {
  image = "example"
  name = "example"
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

- `image` (String, Required) The full image location, including the container registry, image name, and version tag.
- `name` (String, Required) AWX value stored in `name`.

### Optional

- `credential_id` (Number, Optional) Numeric ID of the related AWX credential object.
- `description` (String, Optional) AWX value stored in `description`.
- `organization_id` (Number, Optional) The organization used to determine access to this execution environment.
- `pull` (String, Optional) Pull image before running?
  - `always` - Always pull container before running.
  - `missing` - Only pull the image if not present before running.
  - `never` - Never pull container before running.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_execution_environment.example 42
```

## Further Reading

- [AWX Execution Environments](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/execution_environments.html)
