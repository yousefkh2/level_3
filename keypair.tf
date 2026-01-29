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
