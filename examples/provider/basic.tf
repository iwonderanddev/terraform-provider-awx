terraform {
  required_providers {
    awx = {
      source  = "damien/awx"
      version = "~> 0.1"
    }
  }
}

provider "awx" {
  base_url  = var.awx_base_url
  username  = var.awx_username
  password  = var.awx_password
}
