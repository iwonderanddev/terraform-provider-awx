# Resource: awx_group

Manages AWX `groups` objects.

## Example Usage

```hcl
resource "awx_group" "example" {
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
- `inventory` (Required) Managed field from AWX OpenAPI schema.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `variables` (Optional) Group variables in JSON or YAML format.

## Attributes Reference

- `id` (String) Numeric AWX object identifier.

## Import

```bash
terraform import awx_group.example 42
```
