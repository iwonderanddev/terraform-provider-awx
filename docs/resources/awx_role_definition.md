# Resource: awx_role_definition

Manages AWX `role_definitions` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_role_definition" "example" {
  name = "example"
  permissions = ["awx.view_inventory"]
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

- `name` (String, Required) The name of this role.
- `permissions` (List of String, Required) List of permission strings for `permissions`.

### Optional

- `content_type` (String, Optional) String to use for references to this type from other models in the API.
- `description` (String, Optional) A description of this role.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_role_definition.example 42
```

## Further Reading

- [AWX Role-Based Access Controls](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/rbac.html)
