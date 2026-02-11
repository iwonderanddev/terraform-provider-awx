# Resource: awx_workflow_job_template_node_credential_association

Manages `workflow_job_template_node_credential_association` relationships between `workflow_job_template_nodes` and `credentials` objects.

Breaking change: use `workflow_job_template_node_id` and `credential_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_workflow_job_template_node_credential_association" "example" {
  workflow_job_template_node_id = 12
  credential_id  = 34
}
```

## Argument Reference

- `workflow_job_template_node_id` (Number, Required) Parent object numeric ID.
- `credential_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `workflow_job_template_node_id` (Number) Parent object numeric ID.
- `credential_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_workflow_job_template_node_credential_association.example 12:34
```
