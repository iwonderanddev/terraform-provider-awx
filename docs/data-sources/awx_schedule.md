# Data Source: awx_schedule

Reads AWX `schedules` objects.

## Example Usage

```hcl
data "awx_schedule" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `description` (String, Read-Only) Value for `description`.
- `diff_mode` (Boolean, Read-Only) Controls whether `diff_mode` is enabled in AWX.
- `enabled` (Boolean, Read-Only) Enables processing of this schedule.
- `execution_environment_id` (Number, Read-Only) The container image to be used for execution.
- `extra_data` (Object, Read-Only) Structured extra data as a Terraform object.
- `forks` (Number, Read-Only) Numeric setting for `forks`.
- `inventory_id` (Number, Read-Only) Inventory applied as a prompt, assuming job template prompts for inventory
- `job_slice_count` (Number, Read-Only) Numeric setting for `job_slice_count`.
- `job_tags` (String, Read-Only) Value for `job_tags`.
- `job_type` (String, Read-Only) * `run` - Run
  - `check` - Check
- `limit` (String, Read-Only) Value for `limit`.
- `name` (String, Read-Only) Value for `name`.
- `rrule` (String, Read-Only) A value representing the schedules iCal recurrence rule.
- `scm_branch` (String, Read-Only) Value for `scm_branch`.
- `skip_tags` (String, Read-Only) Value for `skip_tags`.
- `timeout` (Number, Read-Only) Numeric setting for `timeout`.
- `unified_job_template_id` (Number, Read-Only) Numeric ID of the related AWX unified job template object.
- `verbosity` (Number, Read-Only) * `0` - 0 (Normal)
  - `1` - 1 (Verbose)
  - `2` - 2 (More Verbose)
  - `3` - 3 (Debug)
  - `4` - 4 (Connection Debug)
  - `5` - 5 (WinRM Debug)

## Further Reading

- [AWX Schedules](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/scheduling.html)
