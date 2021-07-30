resource "helm_release" "opa" {
  name       = "opa"
  namespace  = kubernetes_namespace.rode.metadata[0].name
  chart      = "opa"
  repository = "https://open-policy-agent.github.io/kube-mgmt/charts"
  version    = "0.12.1"

  set {
    name  = "imageTag"
    value = "0.24.0"
  }

  set {
    name  = "mgmt.enabled"
    value = false
  }
}

resource "kubernetes_deployment" "opa" {
  metadata {
    name = kubernetes_namespace.rode.metadata[0].name
    labels {
      app = "opa"
    }
  }
  spec {
    template {
      metadata {
        labels {
          app = "opa"
        }
      }
      spec {
        container {
          image = "openpolicyagent/opa:0.24.0"
          name  = "opa"
          liveness_probe {
            http_get {
              path = "/health"
              port = 8181
            }
            initial_delay_seconds = 3
            period_seconds        = 5
          }
          readiness_probe {
            http_get {
              path = "/health"
              port = 8181
            }
            initial_delay_seconds = 3
            period_seconds        = 5
          }
          args  = [
            "run",
            "--server"]
        }
      }
    }
  }
}

resource "kubernetes_service" "opa" {
  metadata {
    name      = "opa"
    namespace = kubernetes_namespace.rode.metadata[0].name
  }
  spec {
    selector = {
      app = kubernetes_deployment.opa.metadata[0].labels.app
    }

    port {
      port        = 8181
      target_port = 8181
    }

    type = "ClusterIP"
  }
}
