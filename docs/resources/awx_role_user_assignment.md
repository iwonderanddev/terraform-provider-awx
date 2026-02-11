# Resource: awx_role_user_assignment

Manages AWX `role_user_assignments` objects.

This endpoint does not support in-place updates; Terraform replaces the resource when arguments change.

## Example Usage

```hcl
resource "awx_role_user_assignment" "example" {
  role_definition_id = awx_role_definition.example.id
  related = { key = "value" }
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `content_type` (Optional) String to use for references to this type from other models in the API.
- `created` (Optional) The date/time this resource was created.
- `created_by_id` (Optional) The user who created this resource.
- `id` (Optional) Managed field from AWX OpenAPI schema.
- `object_ansible_id` (Optional) The resource id of the object this role applies to. An alternative to the object_id field.
- `object_id` (Optional) The primary key of the object this assignment applies to; null value indicates system-wide assignment.
- `related` (Optional) Terraform object value. Managed field from AWX OpenAPI schema.
- `role_definition_id` (Required) The role definition which defines permissions conveyed by this assignment.
- `summary_fields` (Optional) Terraform object value. Managed field from AWX OpenAPI schema.
- `url` (Optional) Managed field from AWX OpenAPI schema.
- `user_id` (Optional) Managed field from AWX OpenAPI schema.
- `user_ansible_id` (Optional) The resource ID of the user who will receive permissions from this assignment. An alternative to user field.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_role_user_assignment.example 42
```
