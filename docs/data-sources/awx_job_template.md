# Data Source: awx_job_template

Reads AWX `job_templates` objects.

## Example Usage

```hcl
data "awx_job_template" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `allow_simultaneous` (boolean)
- `ask_credential_on_launch` (boolean)
- `ask_diff_mode_on_launch` (boolean)
- `ask_execution_environment_on_launch` (boolean)
- `ask_forks_on_launch` (boolean)
- `ask_instance_groups_on_launch` (boolean)
- `ask_inventory_on_launch` (boolean)
- `ask_job_slice_count_on_launch` (boolean)
- `ask_job_type_on_launch` (boolean)
- `ask_labels_on_launch` (boolean)
- `ask_limit_on_launch` (boolean)
- `ask_scm_branch_on_launch` (boolean)
- `ask_skip_tags_on_launch` (boolean)
- `ask_tags_on_launch` (boolean)
- `ask_timeout_on_launch` (boolean)
- `ask_variables_on_launch` (boolean)
- `ask_verbosity_on_launch` (boolean)
- `become_enabled` (boolean)
- `description` (string)
- `diff_mode` (boolean)
- `execution_environment` (integer)
- `extra_vars` (string)
- `force_handlers` (boolean)
- `forks` (integer)
- `host_config_key` (string)
- `inventory` (integer)
- `job_slice_count` (integer)
- `job_tags` (string)
- `job_type` (string)
- `limit` (string)
- `name` (string)
- `opa_query_path` (string)
- `playbook` (string)
- `prevent_instance_group_fallback` (boolean)
- `project` (integer)
- `scm_branch` (string)
- `skip_tags` (string)
- `start_at_task` (string)
- `survey_enabled` (boolean)
- `timeout` (integer)
- `use_fact_cache` (boolean)
- `verbosity` (integer)
- `webhook_credential` (integer)
- `webhook_service` (string)
