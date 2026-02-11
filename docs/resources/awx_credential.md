# Resource: awx_credential

Manages AWX `credentials` objects.

## Example Usage

```hcl
resource "awx_credential" "example" {
  credential_type_id = awx_credential_type.example.id
  name = "example"
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `credential_type_id` (Required) Specify the type of credential you want to create. Refer to the documentation for details on each type.
- `description` (Optional) Managed field from AWX OpenAPI schema.
- `inputs` (Optional, Computed, Sensitive) Credential inputs are a write-only sensitive Terraform object payload.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `organization_id` (Optional) Inherit permissions from organization roles. If provided on creation, do not give either user or team.
- `team_id` (Optional, Sensitive) Write-only field used to add team to owner role. If provided, do not give either user or organization. Only valid for creation.
- `user_id` (Optional, Sensitive) Write-only field used to add user to owner role. If provided, do not give either team or organization. Only valid for creation.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_credential.example 42
```
