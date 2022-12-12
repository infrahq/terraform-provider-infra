// Get all users
data "infra_users" "all" {}

output "my_users" {
  value = data.infra_users.all.users
}

// Get `admin@example.com` user
data "infra_users" "admin" {
  filter {
    name = "admin@example.com"
  }
}

// Get all users who belong to the `administrator` group
data "infra_users" "administrators" {
  filter {
    group_name = "Administrators"
  }
}
