# Generate SSH key pair
resource "tls_private_key" "ssh_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

# Create OpenStack keypair
resource "openstack_compute_keypair_v2" "k8s_keypair" {
  name       = "${var.cluster_name}-keypair"
  public_key = tls_private_key.ssh_key.public_key_openssh
}

# Save private key to local file (optional but useful)
resource "local_file" "private_key" {
  content         = tls_private_key.ssh_key.private_key_pem
  filename        = "${path.module}/${var.cluster_name}-key.pem"
  file_permission = "0600"
}

# Automatically configure SSH for out-of-the-box access
resource "null_resource" "ssh_config" {
  depends_on = [
    local_file.private_key,
    openstack_networking_floatingip_v2.fip
  ]

  provisioner "local-exec" {
    command = <<-EOT
      # Create SSH config directory if it doesn't exist
      mkdir -p ~/.ssh
      chmod 700 ~/.ssh
      
      # Backup existing config
      [ -f ~/.ssh/config ] && cp ~/.ssh/config ~/.ssh/config.backup || true
      
      # Remove old entries for this cluster
      sed -i.bak '/# START ${var.cluster_name}/,/# END ${var.cluster_name}/d' ~/.ssh/config 2>/dev/null || true
      
      # Add new SSH config entries
      cat >> ~/.ssh/config << 'EOF'

# START ${var.cluster_name}
%{for k, fip in openstack_networking_floatingip_v2.fip~}
Host ${fip.address} ${var.cluster_name}-${k}
    HostName ${fip.address}
    User debian
    IdentityFile ${abspath(local_file.private_key.filename)}
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
    LogLevel ERROR

%{endfor~}
# END ${var.cluster_name}
EOF
      
      chmod 600 ~/.ssh/config
      echo "✓ SSH configured! You can now run: ssh ${var.cluster_name}-server"
    EOT
  }

  triggers = {
    floating_ips = jsonencode([for k, fip in openstack_networking_floatingip_v2.fip : fip.address])
    key_file     = local_file.private_key.filename
  }
}
