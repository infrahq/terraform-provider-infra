# Create an access key. The access key will expire after 30 days.
resource "infra_access_key" "example" {
  name = "example"
}

# Create an access key that expires in 90 days.
resource "infra_access_key" "expires_in" {
  name       = "expires_in"
  expires_in = "2160h0m0s" # 24h * 90
}

# Create an access key that expires on December 31, 2023 at 23:59:59 UTC.
resource "infra_access_key" "expires_at" {
  name       = "expires_at"
  expires_at = "2023-12-31T23:59:59Z"
}

# Create an access key that expires in 3 days if not used.
resource "infra_access_key" "inactivity_timeout" {
  name               = "inactivity_timeout"
  inactivity_timeout = "72h"
}

# Create an access key for user@example.com.
resource "infra_access_key" "user_email" {
  name       = "user_email"
  user_email = "user@example.com"
}

# Create a connector access key and install the infra connector using Helm
resource "infra_access_key" "connector" {
  connector_access_key = true
}

resource "helm_release" "infra" {
  repository = "https://helm.infrahq.com"
  chart      = "infra"
  version    = "0.20.6"

  name             = "infra-connector"
  namespace        = "infrahq"
  create_namespace = true

  set {
    name  = "connector.config.accessKey"
    value = infra_access_key.connector.secret
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
