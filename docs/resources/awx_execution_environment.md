# Resource: awx_execution_environment

Manages AWX `execution_environments` objects.

## Example Usage

```hcl
resource "awx_execution_environment" "example" {
  image = "example"
  name = "example"
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `credential` (Optional) Managed field from AWX OpenAPI schema.
- `description` (Optional) Managed field from AWX OpenAPI schema.
- `image` (Required) The full image location, including the container registry, image name, and version tag.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `organization` (Optional) The organization used to determine access to this execution environment.
- `pull` (Optional) Pull image before running?
  - `always` - Always pull container before running.
  - `missing` - Only pull the image if not present before running.
  - `never` - Never pull container before running.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_execution_environment.example 42
```
