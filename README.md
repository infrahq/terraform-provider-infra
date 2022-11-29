# Terraform Provider for Infra

This repository is a [Terraform](https://www.terraform.io) provider for managing [Infra](https://www.infrahq.com) resources. This provider is maintained by the Infra development team.

_Terraform Infra provider is under active development._

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.18

## Building the provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command. This will put the provider binary in the `$GOPATH/bin` directory.

```shell
go install
```

## Using the provider

To use a published version of the provider, run `terraform init` to automatically install the latest provider.

To use a not-yet-published version of the provider, download the [latest][1] release from GitHub, unpack the file, and add the binary to `~/.terraform.d/plugins/registry.terraform.io/infrahq/infra/${version}/${target}`. Afterwards, run `terraform init` to initialize the provider.

To use a custom built version of the provider, follow the steps below:

1. Build the provider binary using the steps [above](#building-the-provider).
1. Create a [`~/.terraformrc`](https://developer.hashicorp.com/terraform/cli/config/config-file) with the following content. Ensure the full path of the directory the binary built in the previous step lives in is passed in to `dev_overrides`.

    ```terraform
    provider_installation {

      # Use "$GOPATH/bin" as an overridden package
      # directory for the infrahq/infra provider. This disables the version and
      # checksum verifications for this provider and forces Terraform to look for
      # the infra provider plugin in the given directory.
      dev_overrides {
        "infrahq/infra" = "<full path to $GOPATH/bin>"
      }

      # For all other providers, install them directly from their origin provider
      # registries as normal. If you omit this, Terraform will _only_ use
      # the dev_overrides block, and so no other providers will be available.
      direct {}
    }
    ```

1. Configure the provider in Terraform.

    ```terraform
    terraform {
      required_providers {
        infra = {
          source  = "infrahq/infra"
        }
      }
    }

    provider "infra" {
      access_key = "xxxxxxxxxx.yyyyyyyyyyyyyyyyyyyyyyyy"
    }
    ```

1. Create your first Infra resources! Examples for all resources and data sources are available in `examples/`.

## Developing the provider

To build the provider, follow the steps above.

To update the documentation, run `go generate`.

To run the test suite, run `make testacc`.
