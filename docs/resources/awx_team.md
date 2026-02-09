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

- `description` (Optional) Managed field from AWX OpenAPI schema.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `organization` (Required) Managed field from AWX OpenAPI schema.

## Attributes Reference

- `id` (String) Numeric AWX object identifier.

## Import

```bash
terraform import awx_team.example 42
```
