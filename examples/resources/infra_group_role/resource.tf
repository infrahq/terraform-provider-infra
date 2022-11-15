resource "infra_group" "example" {
  name = "Example"
}

# Assign Infra "admin" to group
resource "infra_group_role" "admin" {
  group_id = infra_group.example.id
  role     = "admin"
}

# Assign Infra "view" to group by name
resource "infra_group_role" "view" {
  group_name = "Example"
  role       = "view"
}
