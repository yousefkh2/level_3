terraform {
  required_version = ">= 1.6.0"
  required_providers {
    stackit = {
      source  = "stackitcloud/stackit"
    }
    local = {
      source = "hashicorp/local"
      version = "2.5.1"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.12"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.25"
    }
  }
}

provider "stackit" {
  default_region                   = "eu01"
  service_account_key_path = "${path.module}/../../sa-key-2fed783a-c9ca-4e49-870b-b65b36e5e728.json"
}

variable "project_id" {
  type        = string
  description = "STACKIT Project ID"
}

variable "cluster_name" {
  type        = string
  default     = "yousef-ske"
}

resource "stackit_ske_cluster" "this" {
  project_id = var.project_id
  name       = var.cluster_name

  node_pools = [
    {
      name                      = "np-1"
      machine_type              = "g1a.2d"
      os_name                   = "ubuntu"
      os_version_min            = "2204.20250728.0"
      minimum                   = 1
      maximum                   = 1
      availability_zones        = ["eu01-1"]

      volume_type = "storage_premium_perf6"
      volume_size = 100
    }
  ]

  extensions = {
    dns = {
      enabled = true
      zones   = []
    }
  }
}

# 1. Generate the Kubeconfig
resource "stackit_ske_kubeconfig" "this" {
  project_id   = var.project_id
  cluster_name = stackit_ske_cluster.this.name
  
  # Optional: Set expiration (defaults to 1h if unset)
  # refresh = true ensures TF updates it if it expires
  refresh = true 
  expiration = 15552000 
}

# 2. Save it to a file on your local machine
resource "local_file" "kubeconfig" {
  content  = stackit_ske_kubeconfig.this.kube_config
  filename = "${path.module}/kubeconfig.yaml"
  file_permission = "0600" # Secure the file so only you can read it
}