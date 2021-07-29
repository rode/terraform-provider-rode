terraform {
  required_providers {
    helm       = {
      source  = "hashicorp/helm"
      version = "2.2.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.3.2"
    }
    random     = {
      source  = "hashicorp/random"
      version = "3.0.1"
    }
  }
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig
  }
}

provider "kubernetes" {
  config_path = var.kubeconfig
}

resource "kubernetes_namespace" "rode" {
  metadata {
    name = "rode"
  }
}

resource "helm_release" "elasticsearch" {
  name       = "elasticsearch"
  namespace  = kubernetes_namespace.rode.metadata[0].name
  chart      = "elasticsearch"
  repository = "https://helm.elastic.co"
  version    = "7.10.1"
  wait       = true

  values = [
<<EOF
clusterHealthCheckParams: wait_for_status=yellow&timeout=1s
esJavaOpts: "-Xmx512m -Xms512m"
persistence:
  enabled: false
replicas: 1
resources:
  requests:
    cpu: 100m
    memory: 512M
sysctlInitContainer:
  enabled: false
tests:
  enabled: false
EOF
]
}

