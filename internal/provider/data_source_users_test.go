package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUsers(t *testing.T) {
	email := randomEmail()
	name := randomName()

	dataSourceName := fmt.Sprintf("data.infra_users.%s", t.Name())

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGroup(t, name),
					testAccResourceGroupMembership(t),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("infra_user.%s", t.Name()), "name", email),
					resource.TestCheckResourceAttr(fmt.Sprintf("infra_group.%s", t.Name()), "name", name),
					resource.TestCheckResourceAttr(fmt.Sprintf("infra_group_membership.%s", t.Name()), "user_name", email),
					resource.TestCheckResourceAttr(fmt.Sprintf("infra_group_membership.%s", t.Name()), "group_name", name),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGroup(t, name),
					testAccResourceGroupMembership(t),
					testAccDataSourceUsers_filterByName(t, email),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "users.0.name", email),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.groups.#", "0"),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGroup(t, name),
					testAccResourceGroupMembership(t),
					testAccDataSourceUsers_filterByGroupNameIncludeGroups(t, name),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "users.0.name", email),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.groups.0", name),
				),
			},
		},
	})
}

func testAccDataSourceUsers_filterByName(t *testing.T, email string) string {
	return fmt.Sprintf(`
data "infra_users" "%[1]s" {
	filter {
		name = "%[2]s"
	}
}`, t.Name(), email)
}

func testAccDataSourceUsers_filterByGroupNameIncludeGroups(t *testing.T, name string) string {
	return fmt.Sprintf(`
data "infra_users" "%[1]s" {
	filter {
		group_name = "%[2]s"
	}

	include_groups = true
}`, t.Name(), name)
}
