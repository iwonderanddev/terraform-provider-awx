# Resource: awx_organization_credential_association

Manages `organization_credential_association` relationships between `organizations` and `credentials` objects.

Breaking change: use `organization_id` and `credential_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_organization_credential_association" "example" {
  organization_id = 12
  credential_id  = 34
}
```

## Argument Reference

- `organization_id` (Number, Required) Parent object numeric ID.
- `credential_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `organization_id` (Number) Parent object numeric ID.
- `credential_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_organization_credential_association.example 12:34
```
