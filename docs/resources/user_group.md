---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "infra_user_group Resource - terraform-provider-infra"
subcategory: ""
description: |-
  Provides an Infra user grant. This resource can be used to assign groups to users.
---

# infra_user_group (Resource)

Provides an Infra user grant. This resource can be used to assign groups to users.

## Example Usage

```terraform
resource "infra_user" "example" {
  email = "example@example.com"
}

resuorce "infra_group" "example" {
  name = "Example"
}

# Assign a user to a group.
resource "infra_user_group" "example" {
  user_id  = infra_user.example.id
  group_id = infra_group.example.id
}

# Assign a user to a group.
resource "infra_user_group" "example" {
  user_email = "example@example.com"
  group_name = "Example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `group_id` (String) The ID of the group to assign to the user. One of `group_id`, `group_name` must be set.
- `group_name` (String) The name of the group to assign to the user. One of `group_id`, `group_name` must be set.
- `user_email` (String) The email of the user to assign to the group. One of `user_id`, `user_email` must be set.
- `user_id` (String) The ID of the user to assign to the group. One of `user_id`, `user_email` must be set.

### Read-Only

- `id` (String) The ID of this resource.

