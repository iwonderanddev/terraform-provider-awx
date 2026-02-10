# Resource: awx_job_template

Manages AWX `job_templates` objects.

## Example Usage

```hcl
resource "awx_job_template" "example" {
  name = "example"
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `allow_simultaneous` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_credential_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_diff_mode_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_execution_environment_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_forks_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_instance_groups_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_inventory_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_job_slice_count_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_job_type_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_labels_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_limit_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_scm_branch_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_skip_tags_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_tags_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_timeout_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_variables_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_verbosity_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `become_enabled` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `description` (Optional) Managed field from AWX OpenAPI schema.
- `diff_mode` (Optional, Computed) If enabled, textual changes made to any templated files on the host are shown in the standard output
- `execution_environment` (Optional) The container image to be used for execution.
- `extra_vars` (Optional) Managed field from AWX OpenAPI schema.
- `force_handlers` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `forks` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `host_config_key` (Optional) Managed field from AWX OpenAPI schema.
- `inventory` (Optional) Managed field from AWX OpenAPI schema.
- `job_slice_count` (Optional, Computed) The number of jobs to slice into at runtime. Will cause the Job Template to launch a workflow if value is greater than 1.
- `job_tags` (Optional) Managed field from AWX OpenAPI schema.
- `job_type` (Optional, Computed) * `run` - Run
  - `check` - Check
- `limit` (Optional) Managed field from AWX OpenAPI schema.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `opa_query_path` (Optional) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.
- `playbook` (Optional) Managed field from AWX OpenAPI schema.
- `prevent_instance_group_fallback` (Optional, Computed) If enabled, the job template will prevent adding any inventory or organization instance groups to the list of preferred instances groups to run on.If this setting is enabled and you provided an empty list, the global instance groups will be applied.
- `project` (Optional) Managed field from AWX OpenAPI schema.
- `scm_branch` (Optional) Branch to use in job run. Project default used if blank. Only allowed if project allow_override field is set to true.
- `skip_tags` (Optional) Managed field from AWX OpenAPI schema.
- `start_at_task` (Optional) Managed field from AWX OpenAPI schema.
- `survey_enabled` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `timeout` (Optional, Computed) The amount of time (in seconds) to run before the task is canceled.
- `use_fact_cache` (Optional, Computed) If enabled, the service will act as an Ansible Fact Cache Plugin; persisting facts at the end of a playbook run to the database and caching facts for use by Ansible.
- `verbosity` (Optional, Computed) * `0` - 0 (Normal)
  - `1` - 1 (Verbose)
  - `2` - 2 (More Verbose)
  - `3` - 3 (Debug)
  - `4` - 4 (Connection Debug)
  - `5` - 5 (WinRM Debug)
- `webhook_credential` (Optional) Personal Access Token for posting back the status to the service API
- `webhook_service` (Optional) Service that webhook requests will be accepted from
  - `github` - GitHub
  - `gitlab` - GitLab
  - `bitbucket_dc` - BitBucket DataCenter

## Attributes Reference

- `id` (String) Numeric AWX object identifier.

## Import

```bash
terraform import awx_job_template.example 42
```
