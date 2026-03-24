# Resource: awx_credential

Manages AWX credentials used by job templates, inventory sources, and integrations.

## Example Usage

### Machine credential

```hcl
resource "awx_organization" "platform" {
  name = "platform"
}

data "awx_credential_type" "machine" {
  name = "Machine"
}

resource "awx_credential" "machine" {
  name               = "linux-ssh"
  credential_type_id = data.awx_credential_type.machine.id
  organization_id    = awx_organization.platform.id
  inputs = {
    username     = "ec2-user"
    ssh_key_data = var.machine_private_key
  }
}
```

### Source control token credential

```hcl
resource "awx_organization" "platform" {
  name = "platform"
}

data "awx_credential_type" "scm" {
  name = "Source Control"
}

resource "awx_credential" "git_token" {
  name               = "git-token"
  credential_type_id = data.awx_credential_type.scm.id
  organization_id    = awx_organization.platform.id
  inputs = {
    username = "git"
    password = var.git_pat
  }
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

- `credential_type_id` (Number, Required) Numeric ID of the credential type definition (for example Machine, Source Control, or Vault).
- `name` (String, Required) Credential name shown in AWX.

### Optional

- `description` (String, Optional) Optional explanation of credential usage.
- `inputs` (Object, Optional, Computed, Sensitive, Write-Only) Object containing credential input fields required by the selected credential type.
- `organization_id` (Number, Optional) Numeric ID of the organization that owns this credential.
- `team_id` (Number, Optional, Sensitive, Write-Only) Numeric ID of a team granted owner access when the credential is created.
- `user_id` (Number, Optional, Sensitive, Write-Only) Numeric ID of a user granted owner access when the credential is created.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_credential.example 42
```

## Further Reading

- [AWX Credentials](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credentials.html)
