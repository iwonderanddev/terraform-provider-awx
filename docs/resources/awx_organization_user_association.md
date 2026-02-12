# Resource: awx_organization_user_association

Manages `organization_user_association` relationships between `organizations`
and `users` objects.

## Example Usage

```hcl
resource "awx_organization_user_association" "example" {
  organization_id = 12
  user_id  = 34
}
```

## Argument Reference

- `organization_id` (Number, Required) Parent object numeric ID.
- `user_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `organization_id` (Number) Parent object numeric ID.
- `user_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_organization_user_association.example \
  12:34
```
