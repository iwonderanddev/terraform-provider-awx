# Resource: awx_workflow_job_node_instance_group_association

Manages `workflow_job_node_instance_group_association` relationships between `workflow_job_nodes` and `instance_groups` objects.

Breaking change: use `workflow_job_node_id` and `instance_group_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_workflow_job_node_instance_group_association" "example" {
  workflow_job_node_id = 12
  instance_group_id  = 34
}
```

## Argument Reference

- `workflow_job_node_id` (Number, Required) Parent object numeric ID.
- `instance_group_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `workflow_job_node_id` (Number) Parent object numeric ID.
- `instance_group_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_workflow_job_node_instance_group_association.example 12:34
```
