terraform {
  required_providers {
    infra = {
      source = "infrahq/infra"
    }
  }
}

# Configure Infra Terraform provider.
provider "infra" {
  access_key = "xxxxxxxxxx.yyyyyyyyyyyyyyyyyyyyyyyy"
}
