# Data Source: awx_job_template

Reads AWX `job_templates` objects.

## Example Usage

```hcl
data "awx_job_template" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `allow_simultaneous` (Boolean, Read-Only) Controls whether `allow_simultaneous` is enabled in AWX.
- `ask_credential_on_launch` (Boolean, Read-Only) Controls whether `ask_credential_on_launch` is enabled in AWX.
- `ask_diff_mode_on_launch` (Boolean, Read-Only) Controls whether `ask_diff_mode_on_launch` is enabled in AWX.
- `ask_execution_environment_on_launch` (Boolean, Read-Only) Controls whether `ask_execution_environment_on_launch` is enabled in AWX.
- `ask_forks_on_launch` (Boolean, Read-Only) Controls whether `ask_forks_on_launch` is enabled in AWX.
- `ask_instance_groups_on_launch` (Boolean, Read-Only) Controls whether `ask_instance_groups_on_launch` is enabled in AWX.
- `ask_inventory_on_launch` (Boolean, Read-Only) When true, AWX prompts for inventory selection during launch.
- `ask_job_slice_count_on_launch` (Boolean, Read-Only) Controls whether `ask_job_slice_count_on_launch` is enabled in AWX.
- `ask_job_type_on_launch` (Boolean, Read-Only) Controls whether `ask_job_type_on_launch` is enabled in AWX.
- `ask_labels_on_launch` (Boolean, Read-Only) Controls whether `ask_labels_on_launch` is enabled in AWX.
- `ask_limit_on_launch` (Boolean, Read-Only) When true, AWX prompts for an Ansible limit pattern at launch.
- `ask_scm_branch_on_launch` (Boolean, Read-Only) Controls whether `ask_scm_branch_on_launch` is enabled in AWX.
- `ask_skip_tags_on_launch` (Boolean, Read-Only) When true, AWX prompts for skip tags at launch.
- `ask_tags_on_launch` (Boolean, Read-Only) When true, AWX prompts for job tags at launch.
- `ask_timeout_on_launch` (Boolean, Read-Only) When true, AWX prompts for a job timeout at launch.
- `ask_variables_on_launch` (Boolean, Read-Only) When true, AWX allows operators to provide extra vars at launch time.
- `ask_verbosity_on_launch` (Boolean, Read-Only) Controls whether `ask_verbosity_on_launch` is enabled in AWX.
- `become_enabled` (Boolean, Read-Only) Controls whether `become_enabled` is enabled in AWX.
- `description` (String, Read-Only) Optional explanation displayed to operators in AWX.
- `diff_mode` (Boolean, Read-Only) If enabled, textual changes made to any templated files on the host are shown in the standard output
- `execution_environment_id` (Number, Read-Only) Numeric ID of the execution environment container image used to run jobs.
- `extra_vars` (Object, Read-Only) Object of default extra variables passed to launched jobs.
- `force_handlers` (Boolean, Read-Only) Controls whether `force_handlers` is enabled in AWX.
- `forks` (Number, Read-Only) Numeric AWX value used for `forks`.
- `host_config_key` (String, Read-Only) AWX value stored in `host_config_key`.
- `inventory_id` (Number, Read-Only) Numeric ID of the default inventory used when the launch does not prompt for inventory.
- `job_slice_count` (Number, Read-Only) The number of jobs to slice into at runtime. Will cause the Job Template to launch a workflow if value is greater than 1.
- `job_tags` (String, Read-Only) AWX value stored in `job_tags`.
- `job_type` (String, Read-Only) Execution mode for the template, usually `run` for normal automation or `check` for dry-run behavior.
- `limit` (String, Read-Only) AWX value stored in `limit`.
- `name` (String, Read-Only) Human-readable name shown in the AWX UI and API listings.
- `opa_query_path` (String, Read-Only) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.
- `playbook` (String, Read-Only) Playbook path inside the project repository, for example `site.yml`.
- `prevent_instance_group_fallback` (Boolean, Read-Only) If enabled, the job template will prevent adding any inventory or organization instance groups to the list of preferred instances groups to run on.If this setting is enabled and you provided an empty list, the global instance groups will be applied.
- `project_id` (Number, Read-Only) Numeric ID of the AWX project that provides the playbook content.
- `scm_branch` (String, Read-Only) Branch to use in job run. Project default used if blank. Only allowed if project allow_override field is set to true.
- `skip_tags` (String, Read-Only) AWX value stored in `skip_tags`.
- `start_at_task` (String, Read-Only) AWX value stored in `start_at_task`.
- `survey_enabled` (Boolean, Read-Only) Enables survey-based prompts for this template.
- `timeout` (Number, Read-Only) The amount of time (in seconds) to run before the task is canceled.
- `use_fact_cache` (Boolean, Read-Only) Enables AWX fact caching for jobs launched from this template.
- `verbosity` (Number, Read-Only) Ansible verbosity level used when running this template.
- `webhook_credential_id` (Number, Read-Only) Numeric ID of the credential used to validate webhook callbacks.
- `webhook_service` (String, Read-Only) Webhook provider AWX accepts for launch requests.

## Further Reading

- [AWX Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/job_templates.html)
