terraform {
  required_providers {
    awx = {
      source = "damien/awx"
    }
  }
}

provider "awx" {
  base_url = "https://awx-dev.iwd.re"
  username = "admin"
  password = "pnCoAFXqXAaA4nleTo8KsMUFYazxnBuI"
}

resource "awx_team" "test" {
  name         = "test"
  organization = 1
}

resource "awx_inventory" "mockshop" {
  name         = "Mockshop"
  organization = awx_organization.mockshop.id
  description  = "All Mockshop servers"
  lifecycle {
    ignore_changes = [
      variables
    ]
  }
}

resource "awx_organization" "mockshop" {
  name        = "Mockshop"
  description = "Mockshop organization"
  # Galaxy credentials id 2 is hardcoded in the AWX instance
  galaxy_credentials = [2]
}
