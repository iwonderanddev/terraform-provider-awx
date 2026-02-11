# Resource: awx_organization_project_association

Manages `organization_project_association` relationships between `organizations` and `projects` objects.

Breaking change: use `organization_id` and `project_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_organization_project_association" "example" {
  organization_id = 12
  project_id  = 34
}
```

## Argument Reference

- `organization_id` (Number, Required) Parent object numeric ID.
- `project_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `organization_id` (Number) Parent object numeric ID.
- `project_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_organization_project_association.example 12:34
```
