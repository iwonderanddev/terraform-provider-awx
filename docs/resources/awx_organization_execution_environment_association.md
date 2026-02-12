# Resource: awx_organization_execution_environment_association

Manages `organization_execution_environment_association` relationships between `organizations`
and `execution_environments` objects.

## Example Usage

```hcl
resource "awx_organization_execution_environment_association" "example" {
  organization_id = 12
  execution_environment_id  = 34
}
```

## Argument Reference

- `organization_id` (Number, Required) Parent object numeric ID.
- `execution_environment_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `organization_id` (Number) Parent object numeric ID.
- `execution_environment_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_organization_execution_environment_association.example \
  12:34
```
