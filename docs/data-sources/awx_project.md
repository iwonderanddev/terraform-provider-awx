# Data Source: awx_project

Reads AWX `projects` objects.

## Example Usage

```hcl
data "awx_project" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `allow_override` (Boolean, Read-Only) Allows job templates to override this project branch at launch.
- `scm_credential_id` (Number, Read-Only) Numeric ID of the source-control credential used to access private repositories.
- `default_environment_id` (Number, Read-Only) Numeric ID of the default execution environment for project jobs.
- `description` (String, Read-Only) Optional explanation displayed to project users.
- `local_path` (String, Read-Only) Local path under PROJECTS_ROOT assigned by AWX. This value is chosen by the server and cannot be set in Terraform.
- `name` (String, Read-Only) Human-readable project name in AWX.
- `organization_id` (Number, Read-Only) Numeric ID of the owning organization.
- `scm_branch` (String, Read-Only) Branch, tag, or revision checked out by AWX.
- `scm_clean` (Boolean, Read-Only) Discard any local changes before syncing the project.
- `scm_delete_on_update` (Boolean, Read-Only) Delete the project before syncing.
- `scm_refspec` (String, Read-Only) For git projects, an additional refspec to fetch.
- `scm_track_submodules` (Boolean, Read-Only) Track submodules latest commits on defined branch.
- `scm_type` (String, Read-Only) Source control backend used by AWX (`git`, `svn`, `archive`, or manual).
- `scm_update_cache_timeout` (Number, Read-Only) Seconds AWX waits before allowing another automatic project update.
- `scm_update_on_launch` (Boolean, Read-Only) When true, AWX updates project content before dependent jobs launch.
- `scm_url` (String, Read-Only) Repository URL or remote archive URL for project content.
- `signature_validation_credential_id` (Number, Read-Only) An optional credential used for validating files in the project against unexpected changes.
- `timeout` (Number, Read-Only) Maximum project update runtime in seconds before AWX cancels the update.

## Further Reading

- [AWX Projects](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/projects.html)
