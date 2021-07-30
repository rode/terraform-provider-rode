resource "helm_release" "grafeas" {
  name       = "grafeas"
  namespace = kubernetes_namespace.rode.metadata[0].name
  chart      = "grafeas-elasticsearch"
  repository = "https://rode.github.io/charts"
  version    = "0.2.0"
  wait       = true
  values = [
  <<EOF
grafeas:
  elasticsearch:
    url: http://elasticsearch-master:9200
    refresh: "true"
    username: "invalid"
    password: "invalid"
elasticsearch:
  enabled: false
image:
  tag: v0.8.2
EOF
  ]

  depends_on = [
    helm_release.elasticsearch,
  ]
}
