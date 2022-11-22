---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "infra_group_grant Resource - terraform-provider-infra"
subcategory: ""
description: |-
  Provides an Infra group grant. This resource can be used to assign grants to groups.
---

# infra_group_grant (Resource)

Provides an Infra group grant. This resource can be used to assign grants to groups.

## Example Usage

```terraform
resource "infra_group" "example" {
  name = "Example"
}

# Grant a group, by ID, the `view` role to a Kubernetes cluster.
resource "infra_group_grant" "view" {
  group_id = infra_group.example.id

  kubernetes {
    cluster = "my_cluster"
    role    = "view"
  }
}

# Grant a group, by name, the `edit` role to a Kubernetes cluster.
resource "infra_group_grant" "edit" {
  group_name = "Example"

  kubernetes {
    cluster = "my_cluster"
    role    = "edit"
  }
}

# Grant a group, by ID, the `admin` role to the `default` namespace in a Kubernetes cluster.
resource "infra_group_grant" "admin" {
  group_id = infra_group.example.id

  kubernetes {
    cluster   = "my_cluster"
    role      = "admin"
    namespace = "default"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `group_id` (String) The ID of the group to assign this grant. One of `group_id`, `group_name` must be set.
- `group_name` (String) The name of the group to assign this grant. One of `group_id`, `group_name` must be set.
- `kubernetes` (Block List, Max: 1) Kubernetes group grant configurations. One of `kubernetes` must be set. (see [below for nested schema](#nestedblock--kubernetes))

### Read-Only

- `id` (String) The grant's unique identifier.

<a id="nestedblock--kubernetes"></a>
### Nested Schema for `kubernetes`

Required:

- `cluster` (String) The name of the Kubernetes cluster to assign to the user.
- `role` (String) The name of the Kubernetes ClusterRole to assign to the group.

Optional:

- `namespace` (String) The namespace of the Kubernetes cluster to assign to the name.

