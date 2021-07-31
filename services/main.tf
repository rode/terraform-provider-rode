terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "2.2.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.3.2"
    }
  }
  required_version = ">= 1.0.0"
}

provider "kubernetes" {
  config_path    = var.kubeconfig
  config_context = var.kubectx
}

provider "helm" {
  kubernetes {
    config_path    = var.kubeconfig
    config_context = var.kubectx
  }
}

resource "kubernetes_namespace" "terraform_provider_rode" {
  metadata {
    name        = "terraform-provider-rode"
    annotations = {}
    labels      = {}
  }
}
