terraform {
  required_version = ">= 0.14.0"
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 3.4.0"
    }
  }
}

provider "openstack" {
  auth_url            = "http://10.0.0.8/identity"
  region              = "RegionOne"
  user_name           = "admin"
  password            = "serr_1"
  tenant_name         = "admin"
  user_domain_id      = "default"
  project_domain_id   = "default"
}