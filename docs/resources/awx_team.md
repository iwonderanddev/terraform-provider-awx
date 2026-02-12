# Resource: awx_team

Manages AWX `teams` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_team" "example" {
  name = "example"
  organization_id = awx_organization.example.id
}
```

## Schema

### Qualifiers

- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Computed`: AWX sets the value during create or refresh.
- `Sensitive`: Terraform redacts the value in normal CLI output.
- `Write-Only`: Sent to AWX during create/update and not read back.

### Required

- `name` (String, Required) Value for `name`.
- `organization_id` (Number, Required) Numeric ID of the related AWX organization object.

### Optional

- `description` (String, Optional) Value for `description`.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_team.example 42
```

## Further Reading

- [AWX Teams](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/teams.html)
