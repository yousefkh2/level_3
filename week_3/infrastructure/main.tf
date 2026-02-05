terraform {
  required_version = ">= 1.6.0"
  required_providers {
    stackit = {
      source  = "stackitcloud/stackit"
    }
  }
}

provider "stackit" {
  default_region = "eu01"
}

variable "cluster_name" {
  type        = string
  default     = "week3-paas"
}

resource "stackit_ske_cluster" "this" {
  project_id = "6f561559-539c-4f64-9615-88f62f68e3ea"
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
}