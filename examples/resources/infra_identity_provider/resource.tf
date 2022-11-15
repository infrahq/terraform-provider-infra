# Connect a generic OIDC identity provider
resource "infra_identity_provider" "example" {
  issuer        = "https://my.oidc.provider.com/"
  client_id     = "example_client_id"
  client_secret = "example_client_secret"
}

# Connect Azure AD
data "azuread_client_config" "current" {}

resource "infra_identity_provider" "azure" {
  client_id     = "example_client_id"
  client_secret = "example_client_secret"

  azure {
    tenant_id = data.azuread_client_config.current.tenant_id
  }
}

# Connect Google without groups
resource "infra_identity_provider" "google" {
  client_id     = "example_client_id"
  client_secret = "example_client_secret"

  google {}
}

# Connect Google with groups
resource "google_service_account" "my_account" {
  account_id = "myaccount"
}

resource "google_service_account_key" "my_key" {
  service_account_id = google_service_account.my_account.name
}

resource "infra_identity_provider" "google_groups" {
  client_id     = "example_client_id"
  client_secret = "example_client_secret"

  google {
    admin_email         = "admin@example.com"
    service_account_key = base64decode(google_service_account_key.my_key.private_key)
  }
}

# Connect Google with groups (manual)
resource "infra_identity_provider" "google_groups" {
  client_id     = "example_client_id"
  client_secret = "example_client_secret"

  google {
    admin_email = "admin@example.com"
    service_account_key = jsonencode({
      private_key : "...",
      client_email : "...",
    })
  }
}

# Connect Okta
data "okta_auth_server" "default" {
  name = "default"
}

resource "okta_app_oauth" "infra" {
  label          = "Infra"
  type           = "web"
  grant_types    = ["authorization_code", "refresh_token"]
  response_types = ["code"]

  redirect_uris = [
    "https://my_organization.infrahq.com/login/callback",
  ]

  groups_claim {
    type        = "FILTER"
    filter_type = "REGEX"
    name        = "groups"
    value       = ".*"
  }

  lifecycle {
    ignore_changes = [
      groups,
    ]
  }
}

data "okta_everyone_group" "everyone" {}

resource "okta_app_group_assignments" "example" {
  app_id = okta_app_oauth.infra.id
  group {
    id = data.okta_everyone_group.everyone.id
  }
}

resource "infra_identity_provider" "okta" {
  client_id     = okta_app_oauth.infra.client_id
  client_secret = okta_app_oauth.infra.client_secret

  okta {
    issuer = data.okta_auth_server.default.issuer
  }
}
