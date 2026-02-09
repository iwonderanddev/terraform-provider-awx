# Resource: awx_project

Manages AWX `projects` objects.

## Example Usage

```hcl
resource "awx_project" "example" {
  name = "example"
}
```

## Argument Reference

- `allow_override` (Optional) Allow changing the SCM branch or revision in a job template that uses this project.
- `credential` (Optional) Managed field from AWX OpenAPI schema.
- `default_environment` (Optional) The default execution environment for jobs run using this project.
- `description` (Optional) Managed field from AWX OpenAPI schema.
- `local_path` (Optional) Local path (relative to PROJECTS_ROOT) containing playbooks and related files for this project.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `organization` (Optional) The organization used to determine access to this template.
- `scm_branch` (Optional) Specific branch, tag or commit to checkout.
- `scm_clean` (Optional) Discard any local changes before syncing the project.
- `scm_delete_on_update` (Optional) Delete the project before syncing.
- `scm_refspec` (Optional) For git projects, an additional refspec to fetch.
- `scm_track_submodules` (Optional) Track submodules latest commits on defined branch.
- `scm_type` (Optional) Specifies the source control system used to store the project.

* `` - Manual
* `git` - Git
* `svn` - Subversion
* `insights` - Red Hat Insights
* `archive` - Remote Archive
- `scm_update_cache_timeout` (Optional) The number of seconds after the last project update ran that a new project update will be launched as a job dependency.
- `scm_update_on_launch` (Optional) Update the project when a job is launched that uses the project.
- `scm_url` (Optional) The location where the project is stored.
- `signature_validation_credential` (Optional) An optional credential used for validating files in the project against unexpected changes.
- `timeout` (Optional) The amount of time (in seconds) to run before the task is canceled.

## Attributes Reference

- `id` (String) Numeric AWX object identifier.

## Import

```bash
terraform import awx_project.example 42
```
