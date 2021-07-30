resource "helm_release" "grafeas" {
  chart      = "grafeas-elasticsearch"
  name       = "grafeas"
  repository = "https://rode.github.io/charts"
  version    = "0.2.0"
  wait       = true
  values = [
  <<EOF
grafeas:
  elasticsearch:
    url: http://elasticsearch-master:9200
    refresh: "true"
elasticsearch:
  enabled: false
image: v0.8.2
EOF
  ]
}