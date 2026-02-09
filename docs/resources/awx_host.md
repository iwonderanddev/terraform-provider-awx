# Resource: awx_host

Manages AWX `hosts` objects.

## Example Usage

```hcl
resource "awx_host" "example" {
  inventory = 1
  name = "example"
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `description` (Optional) Managed field from AWX OpenAPI schema.
- `enabled` (Optional, Computed) Is this host online and available for running jobs?
- `instance_id` (Optional) The value used by the remote inventory source to uniquely identify the host
- `inventory` (Required) Managed field from AWX OpenAPI schema.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `variables` (Optional) Host variables in JSON or YAML format.

## Attributes Reference

- `id` (String) Numeric AWX object identifier.

## Import

```bash
terraform import awx_host.example 42
```
