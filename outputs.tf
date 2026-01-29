


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

output "easy_ssh_commands" {
  description = "Copy/Paste these commands to connect immediately"
  value = {
    for k, v in local.nodes : k => format(
      "ssh -i ./%s-key.pem -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null debian@%s", 
      var.cluster_name, 
      openstack_networking_floatingip_v2.fip[k].address
    )
  }
}