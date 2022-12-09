// Get all destinations
data "infra_destinations" "all" {}

output "my_destinations" {
  value = data.infra_destinations.all.destinations
}
