# Resource: awx_project

Manages AWX projects that supply automation content from source control or archive locations.

## Example Usage

### Git-based project

```hcl
resource "awx_organization" "platform" {
  name = "platform"
}

resource "awx_project" "app" {
  name                 = "app-project"
  organization_id      = awx_organization.platform.id
  scm_type             = "git"
  scm_url              = "https://github.com/example/app-automation.git"
  scm_branch           = "main"
  scm_update_on_launch = true
}
```

### Project with private repository credential

This example assumes an existing source-control credential named `git-deploy-key`.

```hcl
resource "awx_organization" "platform" {
  name = "platform"
}

data "awx_credential" "git_deploy_key" {
  name = "git-deploy-key"
}

resource "awx_project" "private_repo" {
  name                 = "private-automation"
  organization_id      = awx_organization.platform.id
  scm_type             = "git"
  scm_url              = "git@github.com:example/private-automation.git"
  scm_branch           = "release"
  scm_credential_id    = data.awx_credential.git_deploy_key.id
  scm_update_on_launch = true
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

- `name` (String, Required) Human-readable project name in AWX.

### Optional

- `allow_override` (Boolean, Optional, Computed) Allows job templates to override this project branch at launch.
- `scm_credential_id` (Number, Optional) Numeric ID of the source-control credential used to access private repositories.
- `default_environment_id` (Number, Optional) Numeric ID of the default execution environment for project jobs.
- `description` (String, Optional) Optional explanation displayed to project users.
- `organization_id` (Number, Optional) Numeric ID of the owning organization.
- `scm_branch` (String, Optional) Branch, tag, or revision checked out by AWX.
- `scm_clean` (Boolean, Optional, Computed) Discard any local changes before syncing the project.
- `scm_delete_on_update` (Boolean, Optional, Computed) Delete the project before syncing.
- `scm_refspec` (String, Optional) For git projects, an additional refspec to fetch.
- `scm_track_submodules` (Boolean, Optional, Computed) Track submodules latest commits on defined branch.
- `scm_type` (String, Optional) Source control backend used by AWX (`git`, `svn`, `archive`, or manual).
- `scm_update_cache_timeout` (Number, Optional, Computed) Seconds AWX waits before allowing another automatic project update.
- `scm_update_on_launch` (Boolean, Optional, Computed) When true, AWX updates project content before dependent jobs launch.
- `scm_url` (String, Optional) Repository URL or remote archive URL for project content.
- `signature_validation_credential_id` (Number, Optional) An optional credential used for validating files in the project against unexpected changes.
- `timeout` (Number, Optional, Computed) Maximum project update runtime in seconds before AWX cancels the update.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.
- `local_path` (String, Read-Only) Local path under PROJECTS_ROOT assigned by AWX. This value is chosen by the server and cannot be set in Terraform.

## Import

```bash
terraform import awx_project.example 42
```

## Further Reading

- [AWX Projects](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/projects.html)
