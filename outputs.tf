


output "node_fixed_ips" {
  value = {
    for k, port in openstack_networking_port_v2.k8s_ports :
    k => port.all_fixed_ips[0]
  }
}

output "ssh_private_key_path" {
  description = "Path to the generated SSH private key"
  value       = local_file.private_key.filename
}

output "ssh_connection_command" {
  description = "Example SSH connection command"
  value       = "ssh -i ${local_file.private_key.filename} debian@<floating_ip>"
}