resource "awx_organization" "platform" {
  name        = "platform"
  description = "Managed by Terraform"
}

resource "awx_team" "platform" {
  name            = "platform"
  organization_id = awx_organization.platform.id
  description     = "Managed by Terraform"
}

# Object resource imports use numeric AWX IDs.
# terraform import awx_team.platform 42
