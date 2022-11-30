resource "infra_user" "example" {
  name = "example@example.com"
}

# Grant a user, by ID, Infra "admin"
resource "infra_grant" "infra_admin" {
  user_id = infra_user.example.id

  infra {
    role = "admin"
  }
}

# Grant a user, by name, Kubernetes admin
resource "infra_grant" "kubernetes_admin" {
  user_name = "example@example.com"

  kubernetes {
    cluster = "my_cluster"
    role    = "admin"
  }
}

resource "infra_group" "example" {
  name = "Example"
}

# Grant a group, by ID, Kubernetes namespace edit
resource "infra_grant" "kubernetes_namespace_edit" {
  group_id = infra_group.example.id

  kubernetes {
    cluster   = "my_cluster"
    role      = "edit"
    namespace = "default"
  }
}
