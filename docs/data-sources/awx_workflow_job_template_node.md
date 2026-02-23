# Data Source: awx_workflow_job_template_node

Reads AWX `workflow_job_template_nodes` objects.

## Example Usage

```hcl
data "awx_workflow_job_template_node" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `all_parents_must_converge` (Boolean, Read-Only) If enabled then the node will only run if all of the parent nodes have met the criteria to reach this node
- `diff_mode` (Boolean, Read-Only) Controls whether `diff_mode` is enabled in AWX.
- `execution_environment_id` (Number, Read-Only) The container image to be used for execution.
- `extra_data` (Object, Read-Only) Structured extra data as a Terraform object.
- `forks` (Number, Read-Only) Numeric AWX value used for `forks`.
- `identifier` (String, Read-Only) An identifier for this node that is unique within its workflow. It is copied to workflow job nodes corresponding to this node.
- `inventory_id` (Number, Read-Only) Inventory applied as a prompt, assuming job template prompts for inventory
- `job_slice_count` (Number, Read-Only) Numeric AWX value used for `job_slice_count`.
- `job_tags` (String, Read-Only) AWX value stored in `job_tags`.
- `job_type` (String, Read-Only) Allowed values:
  - `run` - Run
  - `check` - Check
- `limit` (String, Read-Only) AWX value stored in `limit`.
- `scm_branch` (String, Read-Only) AWX value stored in `scm_branch`.
- `skip_tags` (String, Read-Only) AWX value stored in `skip_tags`.
- `timeout` (Number, Read-Only) Numeric AWX value used for `timeout`.
- `unified_job_template_id` (Number, Read-Only) Numeric ID of the related AWX unified job template object.
- `verbosity` (Number, Read-Only) Allowed values:
  - `0` - 0 (Normal)
  - `1` - 1 (Verbose)
  - `2` - 2 (More Verbose)
  - `3` - 3 (Debug)
  - `4` - 4 (Connection Debug)
  - `5` - 5 (WinRM Debug)
- `workflow_job_template_id` (Number, Read-Only) Numeric ID of the related AWX workflow job template object.

## Further Reading

- [AWX Workflow Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/workflow_templates.html)
