terraform {
  required_providers {
    awx = {
      source  = "damien/awx"
      version = "~> 0.1"
    }
  }
}

provider "awx" {
  hostname  = var.awx_hostname
  username  = var.awx_username
  password  = var.awx_password
}
