# Helm provider needs kubeconfig
provider "helm" {
  kubernetes {
    config_path = local_file.kubeconfig.filename
  }
}

provider "kubernetes" {
  config_path = local_file.kubeconfig.filename
}

# Install ArgoCD via Helm
resource "helm_release" "argocd" {
  name             = "argocd"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  namespace        = "argocd"
  create_namespace = true
  version          = "7.7.11"

  depends_on = [
    stackit_ske_cluster.this,
    local_file.kubeconfig
  ]
}

# Wait for the CRD to be established + give API discovery a moment
resource "null_resource" "wait_argocd_crd" {
  depends_on = [helm_release.argocd]

  provisioner "local-exec" {
    command = <<EOT
      KUBECONFIG=${local_file.kubeconfig.filename} \
      kubectl wait --for=condition=Established --timeout=300s crd/applications.argoproj.io
    EOT
  }
}

resource "time_sleep" "after_crd" {
  depends_on      = [null_resource.wait_argocd_crd]
  create_duration = "20s"
}

# Apply the root app
# resource "kubernetes_manifest" "argocd_root_app" {
#   manifest = yamldecode(file("${path.module}/../../gitops/root-app.yaml"))
  
#   depends_on = [time_sleep.after_crd]
# }

resource "null_resource" "apply_argocd_root_app" {
  depends_on = [time_sleep.after_crd]

  triggers = {
    # This ensures the command runs again if the file changes
    manifest_sha1 = filesha1("${path.module}/../../gitops/root-app.yaml")
  }

  provisioner "local-exec" {
    command = <<EOT
      KUBECONFIG=${local_file.kubeconfig.filename} \
      kubectl apply -f ${path.module}/../../gitops/root-app.yaml
    EOT
  }
}