data "infra_kubernetes_clusters" "all" {}

output "my_clusters" {
  value = data.infra_kubernetes_clusters.all
}
