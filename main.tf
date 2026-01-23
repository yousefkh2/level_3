terraform {
required_version = ">= 0.14.0"
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.53.0"
    }
  }
}

provider "openstack" {}


########################
# Variables #
########################

variable "instance_name" {
  type    = string
  default = "basic"
}

# variable "key_pair_name" {
#   type = string
# }

variable "network_name" {
  type    = string
  default = "shared"
}


variable "security_groups" {
  type    = list(string)
  default = ["default"]
}


########################
# Lookups (no hardcoded IDs)
########################
data "openstack_images_image_v2" "cirros" {
  name        = "cirros-0.6.3-x86_64-disk"
  most_recent = true
}

data "openstack_compute_flavor_v2" "nano" {
  name = "m1.nano"
}



########################
# Compute Instance
########################
resource "openstack_compute_instance_v2" "basic" {
  name      = var.instance_name
  image_id  = data.openstack_images_image_v2.cirros.id
  flavor_id = data.openstack_compute_flavor_v2.nano.id

  security_groups = var.security_groups

  metadata = {
    this = "that"
  }

  network {
    name = var.network_name
  }
}