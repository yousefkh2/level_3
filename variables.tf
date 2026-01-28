

variable "cluster_name" {
  type    = string
  default = "k8s-lab"
}

variable "network_name" {
  type    = string
  default = "shared"
}

variable "security_groups" {
  type    = list(string)
  default = []
}

variable "key_pair_name" {
  type = string
  description = "Name of existing OpenStack keypair for SSH access"
  # set via tfvars or env var
}

variable "public_network_name" {
  type    = string
  default = "public"
}

variable "image_name" {
  type        = string
  description = "Name of the OS image to use"
  default     = "bookworm"
}

variable "flavor_name" {
  type        = string
  description = "Name of the flavor/instance type"
  default     = "m1.small"
}