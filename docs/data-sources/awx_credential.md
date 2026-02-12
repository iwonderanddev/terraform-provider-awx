# Data Source: awx_credential

Reads AWX `credentials` objects.

## Example Usage

```hcl
data "awx_credential" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `credential_type_id` (Number, Read-Only) Numeric ID of the credential type definition (for example Machine, Source Control, or Vault).
- `description` (String, Read-Only) Optional explanation of credential usage.
- `inputs` (Object, Read-Only, Sensitive) Object containing credential input fields required by the selected credential type.
- `name` (String, Read-Only) Credential name shown in AWX.
- `organization_id` (Number, Read-Only) Numeric ID of the organization that owns this credential.
- `team_id` (Number, Read-Only, Sensitive) Numeric ID of a team granted owner access when the credential is created.
- `user_id` (Number, Read-Only, Sensitive) Numeric ID of a user granted owner access when the credential is created.

## Further Reading

- [AWX Credentials](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credentials.html)
