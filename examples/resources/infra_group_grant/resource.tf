resource "infra_group" "example" {
  name = "Example"
}

# Grant a group, by ID, the `view` role to a Kubernetes cluster.
resource "infra_group_grant" "view" {
  group_id = infra_group.example.id

  kubernetes {
    cluster = "my_cluster"
    role    = "view"
  }
}

# Grant a group, by name, the `edit` role to a Kubernetes cluster.
resource "infra_group_grant" "edit" {
  group_name = "Example"

  kubernetes {
    cluster = "my_cluster"
    role    = "edit"
  }
}

# Grant a group, by ID, the `admin` role to the `default` namespace in a Kubernetes cluster.
resource "infra_group_grant" "admin" {
  group_id = infra_group.example.id

  kubernetes {
    cluster   = "my_cluster"
    role      = "admin"
    namespace = "default"
  }
}
