

data "openstack_compute_flavor_v2" "medium" {
  name = "m1.medium"
}

data "openstack_compute_flavor_v2" "small" {
  name = "m1.small"
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
