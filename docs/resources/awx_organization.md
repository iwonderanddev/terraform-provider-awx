# Resource: awx_organization

Manages AWX `organizations` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_organization" "example" {
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

- None.

### Optional

- `default_environment_id` (Number, Optional) The default execution environment for jobs run by this organization.
- `description` (String, Optional) AWX value stored in `description`.
- `max_hosts` (Number, Optional, Computed) Maximum number of hosts allowed to be managed by this organization.
- `name` (String, Optional) AWX value stored in `name`.
- `opa_query_path` (String, Optional) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_organization.example 42
```

## Further Reading

- [AWX Organizations](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/organizations.html)
