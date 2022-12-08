resource "infra_user" "example" {
  email = "example@example.com"
}

resuorce "infra_group" "example" {
  name = "Example"
}

# Assign a user to a group.
resource "infra_user_group" "example" {
  user_id  = infra_user.example.id
  group_id = infra_group.example.id
}

# Assign a user to a group.
resource "infra_user_group" "example" {
  user_email = "example@example.com"
  group_name = "Example"
}
