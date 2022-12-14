---
page_title: "infra_settings Resource - terraform-provider-infra"
subcategory: ""
description: |-
  Provides Infra organization settings.
  infra_settings behaves differently than normal Terraform resources as settings are
  created with the organization. When a Terraform resource is created, settings automatically
  imported while no action is taken when the resource is deleted.
---

# infra_settings

Provides Infra organization settings.

`infra_settings` behaves differently than normal Terraform resources as settings are
created with the organization. When a Terraform resource is created, settings automatically
imported while no action is taken when the resource is deleted.

## Example Usage

```terraform
resource "infra_settings" "example" {
  password_requirements {
    minimum_length    = 8
    minimum_lowercase = 1
    minimum_uppercase = 1
    minimum_numbers   = 1
    minimum_symbols   = 1
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `password_requirements` (Block List, Max: 1) (see [below for nested schema](#nestedblock--password_requirements))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--password_requirements"></a>
### Nested Schema for `password_requirements`

Optional:

- `minimum_length` (Number) Minimum password length. Default is `8`.
- `minimum_lowercase` (Number) Minimum number of lowercase ASCII letters. Default is `0`.
- `minimum_numbers` (Number) Minimum number of numbers. Default is `0`.
- `minimum_symbols` (Number) Minimum number of symbols. Default is `0`.
- `minimum_uppercase` (Number) Minimum number of uppercase ASCII letters. Default is `0`.

## Import

Import is supported using the following syntax:

```shell
terraform import infra_settings.example example
```