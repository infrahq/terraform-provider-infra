# Configure Infra Terraform provider.
provider "infra" {
  access_key = "xxxxxxxxxx.yyyyyyyyyyyyyyyyyyyyyyyy"
}

# Configure Infra Terraform provider for a custom server.
provider "infra" {
  host       = "https://my.infra.server.com"
  access_key = "xxxxxxxxxx.yyyyyyyyyyyyyyyyyyyyyyyy"
}

# Configure Infra Terraform provider for a custom server
# with a custom trusted server certificate.
provider "infra" {
  host               = "https://my.infra.server.com"
  access_key         = "xxxxxxxxxx.yyyyyyyyyyyyyyyyyyyyyyyy"
  server_certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIE8TCCAtmgAwIBAgIRAJK/+Vzvsa3qzsVf3lOq7nAwDQYJKoZIhvcNAQELBQAw
EDEOMAwGA1UEChMFSW5mcmEwHhcNMjIwNjI5MDEyNjA5WhcNMjMwNjI5MDEzMTA5
WjAQMQ4wDAYDVQQKEwVJbmZyYTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoC
ggIBAM8hcmx6KcQrxKhrEhQldrM81HEdp1/qQfAdMXFKqSKvrR+gGwQO8kPJgkVq
JeuziXUTd0trGt2P0FhfpSwVq+rLuyidDjMcyIXID2HQ6upUpCc1Dsow4RGsXGbU
OQeVNJfCWVTfLtuqScJYQdpzghtVuWMVvsg+UYb722XhlYwGcgIfdeDPatGyRoVN
/CzqC5kl5af0tJIRJ5zPASiQvSIY8+8CGLN8baXP09KTjkHZfMABTa7Htpf71d3H
pk0tLFtzX/P8tPQYDLsJuW/EYYAqN0/YQl+4/wxI5zO5DM+zsErOR6QGiow7sHfh
WYAnY6I8T2DG7aFPa8/K3nRRLNapTgGWuf4gEKK7uPh+WNfYD0+sVaMiE2ZsMswE
aDdM7xlcmjt+/SBMJ9lmk7Ll2vkQiZTE20zMoL/JufTaJlygWZTdLk2T99ccCLBb
+eG56jbmu1xRG+UseqvmaSLTzhOhrVzOCSWw+o0jUfpGO/9DNLkeM++6P9k9vQfo
Sb9qXJCJka+a/CXAAoVjvMHxjFBgY4pzwmgo8K3OiU1BxkWcXFOSMp8cjzT7YhVY
3i1+ZFM4c1MJaTLkufj22yIAd81u+UaW3hMhHGPFsaigGMnazSBFUnAiUvZfgrkX
jJb6bFqlIQJdJIemcVr5cNYD4R0DHNENV+uulyUlwwcfWvE7AgMBAAGjRjBEMA4G
A1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAA
MA8GA1UdEQQIMAaHBH8AAAEwDQYJKoZIhvcNAQELBQADggIBAKQuRDoiScBjC4p7
91k0qYWpraT+zDLvNhxgTnWITje6cb3zZwmALrVfl8FhRLI6ZzjvFWUfVtZCytVO
5Gj5Lo6CflaB8Al9JChqOU4byV5l4YJkeeqiSEO0f/l8R2q8wStKHVMggJYi9ZPp
pLfbVjK/80qQqcxDmbb/ovgKpiA/Ovj5oKOX/uxRbirPhISQx4MuxPzFO6UnGjTe
R6lykdfZLJcC390tmSm0WrxAR7nqE2BxtMQZvN6cnEP0Bw0vOwHbO5bsCgJmyjE6
Fy1abZJwBY5B4RAYneT4BSu4LZXcl11Ow/CKW9WXeroXcOtMgfvzTNDIhIhXJ4RY
WSAEJzFBbaxEy2/PnNzOORhi7ga6B38kv02AsgNY0HyZ6wUCc2gsrPphDYw/A3/W
FvOw73Zj2zqj3H/mDepbFBOFbv+zRlNV12aHY5uq9rRSC9gueP2Lb97PlHXaDsby
6RuPPwGW1vBzmLz2W3gqfMx7xsy5jQJaUeTuV5nI3rUo5ZHniAvlv/eXRenu95u0
1eARbqvpFeNdqOTUnfYpCgEneaK3S+F6xzjnEcjkWN7YdvyWXuVQhPLaRQP5jleI
QLsmSz/7/rCumSogph26nDoTOKxzFS1H5cwwjfToFAKTJX0WaMB0uSuZjbAS55/L
HNDoxYe7Scivi3negOOYzW4bAgp+
-----END CERTIFICATE-----
EOT
}

# Configure Infra Terraform provider for a custom server
# with a custom trusted server certificate file
provider "infra" {
  host                    = "https://my.infra.server.com"
  access_key              = "xxxxxxxxxx.yyyyyyyyyyyyyyyyyyyyyyyy"
  server_certificate_file = "cert.pem"
}
