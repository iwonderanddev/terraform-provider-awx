# Resource: awx_job_template

Manages AWX job templates that launch project playbooks against inventories with optional runtime prompts.

## AWX Concepts

Job templates are reusable launch definitions in AWX. They combine a project, playbook, inventory defaults, and prompt flags so operators can run consistent automation with controlled runtime overrides.

## Example Usage

### Standard project-backed job template

```hcl
resource "awx_organization" "platform" {
  name = "platform"
}

resource "awx_project" "app" {
  name            = "app-project"
  organization_id = awx_organization.platform.id
  scm_type        = "git"
  scm_url         = "https://github.com/example/app-automation.git"
  scm_branch      = "main"
}

resource "awx_inventory" "production" {
  name            = "production"
  organization_id = awx_organization.platform.id
}

resource "awx_job_template" "deploy" {
  name         = "deploy-app"
  job_type     = "run"
  project_id   = awx_project.app.id
  inventory_id = awx_inventory.production.id
  playbook     = "site.yml"
}
```

### Prompted launch with defaults

```hcl
resource "awx_organization" "ops" {
  name = "ops"
}

resource "awx_project" "ops" {
  name            = "ops-automation"
  organization_id = awx_organization.ops.id
  scm_type        = "git"
  scm_url         = "https://github.com/example/ops-automation.git"
}

resource "awx_job_template" "ops_prompted" {
  name                    = "ops-prompted"
  job_type                = "run"
  project_id              = awx_project.ops.id
  playbook                = "operations.yml"
  ask_variables_on_launch = true
  ask_limit_on_launch     = true
  ask_tags_on_launch      = true
  extra_vars = {
    release_channel = "stable"
  }
}
```

### Webhook-enabled job template

This example assumes an existing webhook signing credential named `github-webhook`.

```hcl
resource "awx_organization" "platform" {
  name = "platform"
}

resource "awx_project" "app" {
  name            = "app-project"
  organization_id = awx_organization.platform.id
  scm_type        = "git"
  scm_url         = "https://github.com/example/app-automation.git"
}

resource "awx_inventory" "production" {
  name            = "production"
  organization_id = awx_organization.platform.id
}

data "awx_credential" "github_webhook" {
  name = "github-webhook"
}

resource "awx_job_template" "webhook" {
  name                  = "ci-webhook"
  job_type              = "run"
  project_id            = awx_project.app.id
  inventory_id          = awx_inventory.production.id
  playbook              = "ci.yml"
  webhook_service       = "github"
  webhook_credential_id = data.awx_credential.github_webhook.id
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

- `name` (String, Required) Human-readable name shown in the AWX UI and API listings.

### Optional

- `allow_simultaneous` (Boolean, Optional, Computed) Controls whether `allow_simultaneous` is enabled in AWX.
- `ask_credential_on_launch` (Boolean, Optional, Computed) Controls whether `ask_credential_on_launch` is enabled in AWX.
- `ask_diff_mode_on_launch` (Boolean, Optional, Computed) Controls whether `ask_diff_mode_on_launch` is enabled in AWX.
- `ask_execution_environment_on_launch` (Boolean, Optional, Computed) Controls whether `ask_execution_environment_on_launch` is enabled in AWX.
- `ask_forks_on_launch` (Boolean, Optional, Computed) Controls whether `ask_forks_on_launch` is enabled in AWX.
- `ask_instance_groups_on_launch` (Boolean, Optional, Computed) Controls whether `ask_instance_groups_on_launch` is enabled in AWX.
- `ask_inventory_on_launch` (Boolean, Optional, Computed) When true, AWX prompts for inventory selection during launch.
- `ask_job_slice_count_on_launch` (Boolean, Optional, Computed) Controls whether `ask_job_slice_count_on_launch` is enabled in AWX.
- `ask_job_type_on_launch` (Boolean, Optional, Computed) Controls whether `ask_job_type_on_launch` is enabled in AWX.
- `ask_labels_on_launch` (Boolean, Optional, Computed) Controls whether `ask_labels_on_launch` is enabled in AWX.
- `ask_limit_on_launch` (Boolean, Optional, Computed) When true, AWX prompts for an Ansible limit pattern at launch.
- `ask_scm_branch_on_launch` (Boolean, Optional, Computed) Controls whether `ask_scm_branch_on_launch` is enabled in AWX.
- `ask_skip_tags_on_launch` (Boolean, Optional, Computed) When true, AWX prompts for skip tags at launch.
- `ask_tags_on_launch` (Boolean, Optional, Computed) When true, AWX prompts for job tags at launch.
- `ask_timeout_on_launch` (Boolean, Optional, Computed) When true, AWX prompts for a job timeout at launch.
- `ask_variables_on_launch` (Boolean, Optional, Computed) When true, AWX allows operators to provide extra vars at launch time.
- `ask_verbosity_on_launch` (Boolean, Optional, Computed) Controls whether `ask_verbosity_on_launch` is enabled in AWX.
- `become_enabled` (Boolean, Optional, Computed) Controls whether `become_enabled` is enabled in AWX.
- `description` (String, Optional) Optional explanation displayed to operators in AWX.
- `diff_mode` (Boolean, Optional, Computed) If enabled, textual changes made to any templated files on the host are shown in the standard output
- `execution_environment_id` (Number, Optional) Numeric ID of the execution environment container image used to run jobs.
- `extra_vars` (Object, Optional) Object of default extra variables passed to launched jobs.
- `force_handlers` (Boolean, Optional, Computed) Controls whether `force_handlers` is enabled in AWX.
- `forks` (Number, Optional, Computed) Numeric AWX value used for `forks`.
- `host_config_key` (String, Optional) AWX value stored in `host_config_key`.
- `inventory_id` (Number, Optional) Numeric ID of the default inventory used when the launch does not prompt for inventory.
- `job_slice_count` (Number, Optional, Computed) The number of jobs to slice into at runtime. Will cause the Job Template to launch a workflow if value is greater than 1.
- `job_tags` (String, Optional) AWX value stored in `job_tags`.
- `job_type` (String, Optional, Computed) Execution mode for the template, usually `run` for normal automation or `check` for dry-run behavior.
- `limit` (String, Optional) AWX value stored in `limit`.
- `opa_query_path` (String, Optional) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.
- `playbook` (String, Optional) Playbook path inside the project repository, for example `site.yml`.
- `prevent_instance_group_fallback` (Boolean, Optional, Computed) If enabled, the job template will prevent adding any inventory or organization instance groups to the list of preferred instances groups to run on.If this setting is enabled and you provided an empty list, the global instance groups will be applied.
- `project_id` (Number, Optional) Numeric ID of the AWX project that provides the playbook content.
- `scm_branch` (String, Optional) Branch to use in job run. Project default used if blank. Only allowed if project allow_override field is set to true.
- `skip_tags` (String, Optional) AWX value stored in `skip_tags`.
- `start_at_task` (String, Optional) AWX value stored in `start_at_task`.
- `survey_enabled` (Boolean, Optional, Computed) Enables survey-based prompts for this template.
- `timeout` (Number, Optional, Computed) The amount of time (in seconds) to run before the task is canceled.
- `use_fact_cache` (Boolean, Optional, Computed) Enables AWX fact caching for jobs launched from this template.
- `verbosity` (Number, Optional, Computed) Ansible verbosity level used when running this template.
- `webhook_credential_id` (Number, Optional) Numeric ID of the credential used to validate webhook callbacks.
- `webhook_service` (String, Optional) Webhook provider AWX accepts for launch requests.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_job_template.example 42
```

## Further Reading

- [AWX Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/job_templates.html)
