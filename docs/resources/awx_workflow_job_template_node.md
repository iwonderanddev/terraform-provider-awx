# Resource: awx_workflow_job_template_node

Manages AWX `workflow_job_template_nodes` objects.

## Example Usage

```hcl
resource "awx_workflow_job_template_node" "example" {
  workflow_job_template = 1
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `all_parents_must_converge` (Optional, Computed) If enabled then the node will only run if all of the parent nodes have met the criteria to reach this node
- `diff_mode` (Optional) Managed field from AWX OpenAPI schema.
- `execution_environment` (Optional) The container image to be used for execution.
- `extra_data` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `forks` (Optional) Managed field from AWX OpenAPI schema.
- `identifier` (Optional) An identifier for this node that is unique within its workflow. It is copied to workflow job nodes corresponding to this node.
- `inventory` (Optional) Inventory applied as a prompt, assuming job template prompts for inventory
- `job_slice_count` (Optional) Managed field from AWX OpenAPI schema.
- `job_tags` (Optional) Managed field from AWX OpenAPI schema.
- `job_type` (Optional) * `run` - Run
  - `check` - Check
- `limit` (Optional) Managed field from AWX OpenAPI schema.
- `scm_branch` (Optional) Managed field from AWX OpenAPI schema.
- `skip_tags` (Optional) Managed field from AWX OpenAPI schema.
- `timeout` (Optional) Managed field from AWX OpenAPI schema.
- `unified_job_template` (Optional) Managed field from AWX OpenAPI schema.
- `verbosity` (Optional) * `0` - 0 (Normal)
  - `1` - 1 (Verbose)
  - `2` - 2 (More Verbose)
  - `3` - 3 (Debug)
  - `4` - 4 (Connection Debug)
  - `5` - 5 (WinRM Debug)
- `workflow_job_template` (Required) Managed field from AWX OpenAPI schema.

## Attributes Reference

- `id` (String) Numeric AWX object identifier.

## Import

```bash
terraform import awx_workflow_job_template_node.example 42
```
