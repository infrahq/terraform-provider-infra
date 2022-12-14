---
page_title: "infra_user Resource - terraform-provider-infra"
subcategory: ""
description: |-
  Infra user resource creates a user with a specified name. The name must be an email address.
---

# infra_user

Infra user resource creates a user with a specified name. The name must be an email address.

## Example Usage

```terraform
# Create a user.
resource "infra_user" "example" {
  name = "example@example.com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The user's email address, e.g. `alice@example.com`.

### Optional

- `password` (String, Sensitive) The user's password. This password is one-time use and must be changed before the account can be used. If omitted, a password will be randomly generated. Note: this field will be empty for an imported user.

### Read-Only

- `id` (String) The user's unique identifier.

## Import

Import is supported using the following syntax:

```shell
terraform import infra_user.example <user_id>
```
