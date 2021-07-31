resource "helm_release" "rode" {
  name       = "rode"
  namespace  = kubernetes_namespace.terraform_provider_rode.metadata[0].name
  chart      = "rode"
  repository = "https://rode.github.io/charts"
  version    = "0.3.0"
  wait       = true
  values = [
    <<-EOF
    grafeas-elasticsearch:
      enabled: false
    rode-ui:
      enabled: false
    grafeas:
      host: grafeas-server:8080
    image:
      tag: v${var.rode_version}
    EOF
  ]

  depends_on = [
    helm_release.grafeas,
  ]
}