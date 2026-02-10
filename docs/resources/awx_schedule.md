# Resource: awx_schedule

Manages AWX `schedules` objects.

## Example Usage

```hcl
resource "awx_schedule" "example" {
  name = "example"
  rrule = "example"
  unified_job_template = 1
  extra_data = { key = "value" }
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `description` (Optional) Managed field from AWX OpenAPI schema.
- `diff_mode` (Optional) Managed field from AWX OpenAPI schema.
- `enabled` (Optional, Computed) Enables processing of this schedule.
- `execution_environment` (Optional) The container image to be used for execution.
- `extra_data` (Optional, Computed) Structured extra data as a Terraform object.
- `forks` (Optional) Managed field from AWX OpenAPI schema.
- `inventory` (Optional) Inventory applied as a prompt, assuming job template prompts for inventory
- `job_slice_count` (Optional) Managed field from AWX OpenAPI schema.
- `job_tags` (Optional) Managed field from AWX OpenAPI schema.
- `job_type` (Optional) * `run` - Run
  - `check` - Check
- `limit` (Optional) Managed field from AWX OpenAPI schema.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `rrule` (Required) A value representing the schedules iCal recurrence rule.
- `scm_branch` (Optional) Managed field from AWX OpenAPI schema.
- `skip_tags` (Optional) Managed field from AWX OpenAPI schema.
- `timeout` (Optional) Managed field from AWX OpenAPI schema.
- `unified_job_template` (Required) Managed field from AWX OpenAPI schema.
- `verbosity` (Optional) * `0` - 0 (Normal)
  - `1` - 1 (Verbose)
  - `2` - 2 (More Verbose)
  - `3` - 3 (Debug)
  - `4` - 4 (Connection Debug)
  - `5` - 5 (WinRM Debug)

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_schedule.example 42
```
