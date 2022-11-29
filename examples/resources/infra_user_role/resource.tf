resource "infra_user" "example" {
  email = "example@example.com"
}

# Assign Infra "admin" to user
resource "infra_user_role" "admin" {
  user_id = infra_user.example.id
  role    = "admin"
}

# Assign Infra "view" to user by name
resource "infra_user_role" "view" {
  user_email = "example@example.com"
  role       = "view"
}
