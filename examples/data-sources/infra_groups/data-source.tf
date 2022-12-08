// Get all groups
data "infra_groups" "all" {}

output "my_groups" {
  value = data.infra_groups.all.groups
}

// Get the `engineering` group
data "infra_groups" "engineering" {
  filter {
    name = "Engineering"
  }
}

// Get all groups where `admin@example.com` is a member
data "infra_groups" "" {
  filter {
    user_name = "admin@example.com"
  }
}
