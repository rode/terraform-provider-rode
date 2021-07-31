resource "helm_release" "elasticsearch" {
  name       = "elasticsearch"
  namespace  = kubernetes_namespace.terraform_provider_rode.metadata[0].name
  chart      = "elasticsearch"
  repository = "https://helm.elastic.co"
  version    = "7.10.1"
  wait       = true

  values = [
    <<-EOF
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
