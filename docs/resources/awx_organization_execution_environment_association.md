# Resource: awx_organization_execution_environment_association

Manages AWX associations between `organizations` and `execution_environments` objects.

## Example Usage

```hcl
resource "awx_organization_execution_environment_association" "example" {
  organization_id = 12
  execution_environment_id = 34
}
```

## Schema

### Required

- `organization_id` (Number, Required) Parent object numeric ID.
- `execution_environment_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `organization_id` (Number, Read-Only) Parent object numeric ID.
- `execution_environment_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_organization_execution_environment_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Organizations](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/organizations.html)
- [AWX Execution Environments](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/execution_environments.html)
