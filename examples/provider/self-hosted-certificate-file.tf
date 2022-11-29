# Configure Infra Terraform provider for a custom server
# with a custom trusted server certificate file
provider "infra" {
  host                    = "https://my.infra.server.com"
  access_key              = "xxxxxxxxxx.yyyyyyyyyyyyyyyyyyyyyyyy"
  server_certificate_file = "cert.pem"
}
