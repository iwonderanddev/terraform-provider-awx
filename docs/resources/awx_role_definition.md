# Resource: awx_role_definition

Manages AWX `role_definitions` objects.

## Example Usage

```hcl
resource "awx_role_definition" "example" {
  name = "example"
  permissions = jsonencode(["value"])
}
```

## Argument Reference

- `content_type` (Optional) String to use for references to this type from other models in the API.
- `description` (Optional) A description of this role.
- `name` (Required) The name of this role.
- `permissions` (Required) Managed field from AWX OpenAPI schema.

## Attributes Reference

- `id` (String) Numeric AWX object identifier.

## Import

```bash
terraform import awx_role_definition.example 42
```
