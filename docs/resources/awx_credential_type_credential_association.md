# Resource: awx_credential_type_credential_association

Manages AWX associations between `credential_types` and `credentials` objects.

## Example Usage

```hcl
resource "awx_credential_type_credential_association" "example" {
  credential_type_id = 12
  credential_id = 34
}
```

## Schema

### Required

- `credential_type_id` (Number, Required) Parent object numeric ID.
- `credential_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `credential_type_id` (Number, Read-Only) Parent object numeric ID.
- `credential_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_credential_type_credential_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Credential Types](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credential_types.html)
- [AWX Credentials](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credentials.html)
