# Resource: awx_workflow_job_template_node

Manages AWX `workflow_job_template_nodes` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_workflow_job_template_node" "example" {
  workflow_job_template_id = awx_workflow_job_template.example.id
  extra_data = { key = "value" }
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

- `workflow_job_template_id` (Number, Required) Numeric ID of the related AWX workflow job template object.

### Optional

- `all_parents_must_converge` (Boolean, Optional, Computed) If enabled then the node will only run if all of the parent nodes have met the criteria to reach this node
- `diff_mode` (Boolean, Optional) Controls whether `diff_mode` is enabled in AWX.
- `execution_environment_id` (Number, Optional) The container image to be used for execution.
- `extra_data` (Object, Optional, Computed) Structured extra data as a Terraform object.
- `forks` (Number, Optional) Numeric setting for `forks`.
- `identifier` (String, Optional) An identifier for this node that is unique within its workflow. It is copied to workflow job nodes corresponding to this node.
- `inventory_id` (Number, Optional) Inventory applied as a prompt, assuming job template prompts for inventory
- `job_slice_count` (Number, Optional) Numeric setting for `job_slice_count`.
- `job_tags` (String, Optional) Value for `job_tags`.
- `job_type` (String, Optional) * `run` - Run
  - `check` - Check
- `limit` (String, Optional) Value for `limit`.
- `scm_branch` (String, Optional) Value for `scm_branch`.
- `skip_tags` (String, Optional) Value for `skip_tags`.
- `timeout` (Number, Optional) Numeric setting for `timeout`.
- `unified_job_template_id` (Number, Optional) Numeric ID of the related AWX unified job template object.
- `verbosity` (Number, Optional) * `0` - 0 (Normal)
  - `1` - 1 (Verbose)
  - `2` - 2 (More Verbose)
  - `3` - 3 (Debug)
  - `4` - 4 (Connection Debug)
  - `5` - 5 (WinRM Debug)

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_workflow_job_template_node.example 42
```

## Further Reading

- [AWX Workflow Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/workflow_templates.html)
