resource "awx_team" "platform" {
  name         = "platform"
  organization = 1
  description  = "Managed by Terraform"
}

# Object resource imports use numeric AWX IDs.
# terraform import awx_team.platform 42
