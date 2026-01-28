

data "openstack_images_image_v2" "debian12" {
  name        = var.image_name
  most_recent = true
}

data "openstack_compute_flavor_v2" "medium" {
  name = var.flavor_name
}

data "openstack_networking_network_v2" "k8s_network" {
  name = var.network_name
}

data "openstack_networking_secgroup_v2" "default" {
  name      = "default"
  tenant_id = data.openstack_identity_auth_scope_v3.current.project_id
}

data "openstack_identity_auth_scope_v3" "current" {
  name = "current"
}
