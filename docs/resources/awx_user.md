# Resource: awx_user

Manages AWX `users` objects.

## Example Usage

```hcl
resource "awx_user" "example" {
  username = "example"
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `email` (Optional) Managed field from AWX OpenAPI schema.
- `first_name` (Optional) Managed field from AWX OpenAPI schema.
- `is_superuser` (Optional, Computed) Designates that this user has all permissions without explicitly assigning them.
- `is_system_auditor` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `last_name` (Optional) Managed field from AWX OpenAPI schema.
- `password` (Optional, Sensitive) Field used to change the password.
- `username` (Required) Required. 150 characters or fewer. Letters, digits and @/./+/-/_ only.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_user.example 42
```
