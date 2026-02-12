# Resource: awx_workflow_job_template_node_instance_group_association

Manages AWX associations between `workflow_job_template_nodes` and `instance_groups` objects.

## Example Usage

```hcl
resource "awx_workflow_job_template_node_instance_group_association" "example" {
  workflow_job_template_node_id = 12
  instance_group_id = 34
}
```

## Schema

### Required

- `workflow_job_template_node_id` (Number, Required) Parent object numeric ID.
- `instance_group_id` (Number, Required) Child object numeric ID.

### Read-Only

- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.
- `workflow_job_template_node_id` (Number, Read-Only) Parent object numeric ID.
- `instance_group_id` (Number, Read-Only) Child object numeric ID.
## Import

```bash
terraform import awx_workflow_job_template_node_instance_group_association.example <primary_id>:<related_id>
```

## Further Reading

- [AWX Workflow Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/workflow_templates.html)
- [AWX Instance Groups](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/instance_groups.html)
