locals {
  nodes = {
    server = { name_suffix = "server" }
    node-0 = { name_suffix = "node-0" }
    node-1 = { name_suffix = "node-1" }
  }
}

resource "openstack_compute_instance_v2" "k8s_nodes" {
  for_each = local.nodes

  name      = "${var.cluster_name}-${each.value.name_suffix}"
  image_id  = data.openstack_images_image_v2.debian12.id
  flavor_id = data.openstack_compute_flavor_v2.medium.id

  key_pair = var.key_pair_name

  metadata = {
    role = each.key
  }

  network {
    # Attach the explicit port
    port = openstack_networking_port_v2.k8s_ports[each.key].id
  }
}

# 1. Request a Floating IP for each node
resource "openstack_networking_floatingip_v2" "fip" {
  for_each = local.nodes
  pool     = var.public_network_name
}

# 2. Attach the Floating IP to the Instance
# Using the networking resource 
resource "openstack_networking_floatingip_associate_v2" "fip_assoc" {
  for_each = local.nodes

  # FIX: Use .address (the IP string), NOT .id (the UUID)
  floating_ip = openstack_networking_floatingip_v2.fip[each.key].address

  port_id = openstack_networking_port_v2.k8s_ports[each.key].id
}