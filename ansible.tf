

resource "local_file" "ansible_inventory" {
  # Ensure FIPs are attached before writing the file
  depends_on = [openstack_networking_floatingip_associate_v2.fip_assoc] 
  
  content = <<EOT
[controlplane]
${openstack_compute_instance_v2.k8s_nodes["server"].name} ansible_host=${openstack_networking_floatingip_v2.fip["server"].address} ansible_user=debian ansible_ssh_common_args='-o StrictHostKeyChecking=no'

[workers]
%{ for key, node in openstack_compute_instance_v2.k8s_nodes ~}
%{ if key != "server" ~}
${node.name} ansible_host=${openstack_networking_floatingip_v2.fip[key].address} ansible_user=debian ansible_ssh_common_args='-o StrictHostKeyChecking=no'
%{ endif ~}
%{ endfor ~}

[k8s:children]
controlplane
workers
EOT
  filename = "${path.module}/inventory.ini"
}
