

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

variable "public_network_name" {
  type    = string
  default = "public"
}

variable "flavor_name" {
  type        = string
  description = "Name of the flavor/instance type"
  default     = "m1.medium"
}