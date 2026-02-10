# Resource: awx_team

Manages AWX `teams` objects.

## Example Usage

```hcl
resource "awx_team" "example" {
  name = "example"
  organization = 1
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `description` (Optional) Managed field from AWX OpenAPI schema.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `organization` (Required) Managed field from AWX OpenAPI schema.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_team.example 42
```
