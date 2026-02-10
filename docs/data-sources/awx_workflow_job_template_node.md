# Data Source: awx_workflow_job_template_node

Reads AWX `workflow_job_template_nodes` objects.

## Example Usage

```hcl
data "awx_workflow_job_template_node" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `all_parents_must_converge` (boolean)
- `diff_mode` (boolean)
- `execution_environment` (integer)
- `extra_data` (string)
- `forks` (integer)
- `identifier` (string)
- `inventory` (integer)
- `job_slice_count` (integer)
- `job_tags` (string)
- `job_type` (string)
- `limit` (string)
- `scm_branch` (string)
- `skip_tags` (string)
- `timeout` (integer)
- `unified_job_template` (integer)
- `verbosity` (integer)
- `workflow_job_template` (integer)
