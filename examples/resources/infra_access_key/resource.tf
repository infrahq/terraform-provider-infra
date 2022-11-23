# Create a connector access key that expires
# after 10 years or 7 days without any activity
resource "infra_access_key" "example" {
  expires_in         = "87660h0m0s"
  inactivity_timeout = "168h0m0s"
}

resource "helm_release" "infra" {
  repository = "https://helm.infrahq.com"
  chart      = "infra"
  version    = "0.20.9"

  name             = "infra-connector"
  namespace        = "infrahq"
  create_namespace = true

  set {
    name  = "connector.config.accessKey"
    value = infra_access_key.example.secret
  }

  set {
    name  = "connector.config.server"
    value = "https://api.infrahq.com"
  }

  set {
    name  = "connector.config.name"
    value = "my_cluster"
  }
}
