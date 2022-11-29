resource "infra_user" "example" {
  email = "example@example.com"
}

# Grant a user, by ID, the `view` role to a Kubernetes cluster.
resource "infra_user_grant" "view" {
  user_id = infra_user.example.id

  kubernetes {
    cluster = "my_cluster"
    role    = "view"
  }
}

# Grant a user, by email, the `edit` role to a Kubernetes cluster.
resource "infra_user_grant" "edit" {
  user_email = "example@example.com"

  kubernetes {
    cluster = "my_cluster"
    role    = "edit"
  }
}

# Grant a user, by ID, the `admin` role to the `default` namespace in a Kubernetes cluster.
resource "infra_user_grant" "admin" {
  user_id = infra_user.example.id

  kubernetes {
    cluster   = "my_cluster"
    role      = "admin"
    namespace = "default"
  }
}
