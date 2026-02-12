# Resource: awx_credential_type_credential_association

Manages `credential_type_credential_association` relationships between `credential_types`
and `credentials` objects.

## Example Usage

```hcl
resource "awx_credential_type_credential_association" "example" {
  credential_type_id = 12
  credential_id  = 34
}
```

## Argument Reference

- `credential_type_id` (Number, Required) Parent object numeric ID.
- `credential_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `credential_type_id` (Number) Parent object numeric ID.
- `credential_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_credential_type_credential_association.example \
  12:34
```
