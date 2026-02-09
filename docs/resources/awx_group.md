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
