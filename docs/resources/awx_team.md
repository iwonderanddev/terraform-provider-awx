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
- `Read-Only`: Cannot be set in configuration; Terraform records the value AWX returns.
- `Sensitive`: Terraform redacts the value in normal CLI output.
- `Write-Only`: Sent to AWX during create/update and not read back.

### Required

- `name` (String, Required) AWX value stored in `name`.
- `organization_id` (Number, Required) Numeric ID of the related AWX organization object.

### Optional

- `description` (String, Optional) AWX value stored in `description`.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_team.example 42
```

## Further Reading

- [AWX Teams](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/teams.html)
