# Data Source: awx_schedule

Reads AWX `schedules` objects.

## Example Usage

```hcl
data "awx_schedule" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `description` (string)
- `diff_mode` (boolean)
- `enabled` (boolean)
- `execution_environment` (integer)
- `extra_data` (object)
- `forks` (integer)
- `inventory` (integer)
- `job_slice_count` (integer)
- `job_tags` (string)
- `job_type` (string)
- `limit` (string)
- `name` (string)
- `rrule` (string)
- `scm_branch` (string)
- `skip_tags` (string)
- `timeout` (integer)
- `unified_job_template` (integer)
- `verbosity` (integer)
