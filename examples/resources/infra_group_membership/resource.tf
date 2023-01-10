resource "infra_user" "example" {
  name = "example@example.com"
}

resource "infra_group" "example" {
  name = "Example"
}

# Assign a user to a group.
resource "infra_group_membership" "example" {
  user_id  = infra_user.example.id
  group_id = infra_group.example.id
}

# Assign a user to a group.
resource "infra_group_membership" "example" {
  user_name  = "example@example.com"
  group_name = "Example"
}
