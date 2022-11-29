data "infra_destinations" "all" {}

output "my_clusters" {
  value = data.infra_destinations.all
}
