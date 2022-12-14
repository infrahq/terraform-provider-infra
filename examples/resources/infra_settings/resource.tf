resource "infra_settings" "example" {
  password_requirements {
    minimum_length    = 8
    minimum_lowercase = 1
    minimum_uppercase = 1
    minimum_numbers   = 1
    minimum_symbols   = 1
  }
}
