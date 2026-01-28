#SECURITY GROUPS
# Create a group specifically for the cluster
resource "openstack_networking_secgroup_v2" "k8s_internal" {
  name        = "${var.cluster_name}-internal"
  description = "Allow all internal traffic between k8s nodes"
}

# Allow SSH from anywhere (or restrict to your jumpbox IP if you know it)
resource "openstack_networking_secgroup_v2" "k8s_ssh" {
  name        = "${var.cluster_name}-ssh"
  description = "Allow SSH ingress"
}


# SECURITY RULES
# Allow ALL traffic if it comes from a node inside this same group
resource "openstack_networking_secgroup_rule_v2" "k8s_internal_ipv4" {
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_group_id   = openstack_networking_secgroup_v2.k8s_internal.id
  security_group_id = openstack_networking_secgroup_v2.k8s_internal.id
  # we omit protocol and port ranges to allow ALL
}


resource "openstack_networking_secgroup_rule_v2" "ssh_rule" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.k8s_ssh.id
}


resource "openstack_networking_port_v2" "k8s_ports" {
  for_each       = local.nodes
  name           = "${var.cluster_name}-${each.value.name_suffix}-port"
  network_id     = data.openstack_networking_network_v2.k8s_network.id
  admin_state_up = true

  # We concat:
  # 1. Any extra IDs provided in the variable
  # 2. The LOOKED UP ID of the default group
  # 3. The IDs of the groups we created
  security_group_ids = concat(var.security_groups, [
    data.openstack_networking_secgroup_v2.default.id, 
    openstack_networking_secgroup_v2.k8s_internal.id,
    openstack_networking_secgroup_v2.k8s_ssh.id
  ])
}