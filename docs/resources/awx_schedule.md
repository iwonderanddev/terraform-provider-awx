# Resource: awx_schedule

Manages AWX `schedules` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_schedule" "example" {
  name = "example"
  rrule = "example"
  unified_job_template_id = 1
  extra_data = { key = "value" }
}
```

## Schema

### Qualifiers

- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Computed`: AWX sets the value during create or refresh.
- `Read-Only`: Cannot be set in configuration; Terraform records the value AWX returns.
- `Sensitive`: Terraform redacts the value in normal CLI output.
- `Write-Only`: Sent to AWX during create/update and not read back.

### Required

- `name` (String, Required) AWX value stored in `name`.
- `rrule` (String, Required) A value representing the schedules iCal recurrence rule.
- `unified_job_template_id` (Number, Required) Numeric ID of the related AWX unified job template object.

### Optional

- `description` (String, Optional) AWX value stored in `description`.
- `diff_mode` (Boolean, Optional) Controls whether `diff_mode` is enabled in AWX.
- `enabled` (Boolean, Optional, Computed) Enables processing of this schedule.
- `execution_environment_id` (Number, Optional) The container image to be used for execution.
- `extra_data` (Object, Optional, Computed) Structured extra data as a Terraform object.
- `forks` (Number, Optional) Numeric AWX value used for `forks`.
- `inventory_id` (Number, Optional) Inventory applied as a prompt, assuming job template prompts for inventory
- `job_slice_count` (Number, Optional) Numeric AWX value used for `job_slice_count`.
- `job_tags` (String, Optional) AWX value stored in `job_tags`.
- `job_type` (String, Optional) Allowed values:
  - `run` - Run
  - `check` - Check
- `limit` (String, Optional) AWX value stored in `limit`.
- `scm_branch` (String, Optional) AWX value stored in `scm_branch`.
- `skip_tags` (String, Optional) AWX value stored in `skip_tags`.
- `timeout` (Number, Optional) Numeric AWX value used for `timeout`.
- `verbosity` (Number, Optional) Allowed values:
  - `0` - 0 (Normal)
  - `1` - 1 (Verbose)
  - `2` - 2 (More Verbose)
  - `3` - 3 (Debug)
  - `4` - 4 (Connection Debug)
  - `5` - 5 (WinRM Debug)

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_schedule.example 42
```

## Further Reading

- [AWX Schedules](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/scheduling.html)
