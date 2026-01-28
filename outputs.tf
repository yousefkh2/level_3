


output "node_fixed_ips" {
  value = {
    for k, port in openstack_networking_port_v2.k8s_ports :
    k => port.all_fixed_ips[0]
  }
}