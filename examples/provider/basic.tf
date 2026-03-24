terraform {
  required_providers {
    awx = {
      source  = "iwonderanddev/awx"
      version = "~> 0.2"
    }
  }
}

provider "awx" {
  hostname = var.awx_hostname
  username = var.awx_username
  password = var.awx_password
}
