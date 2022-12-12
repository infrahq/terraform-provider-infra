package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGroups(t *testing.T) {
	email := randomEmail()
	name := randomName()

	dataSourceName := fmt.Sprintf("data.infra_groups.%s", t.Name())

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
					testAccDataSourceGroups_filterByName(t, name),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "groups.0.name", name),
					resource.TestCheckResourceAttr(dataSourceName, "groups.0.users.#", "0"),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGroup(t, name),
					testAccResourceGroupMembership(t),
					testAccDataSourceGroups_filterByUserNameIncludeUsers(t, email),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "groups.0.name", name),
					resource.TestCheckResourceAttr(dataSourceName, "groups.0.users.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "groups.0.users.0", email),
				),
			},
		},
	})
}

func testAccDataSourceGroups_filterByName(t *testing.T, email string) string {
	return fmt.Sprintf(`
data "infra_groups" "%[1]s" {
	filter {
		name = "%[2]s"
	}
}`, t.Name(), email)
}

func testAccDataSourceGroups_filterByUserNameIncludeUsers(t *testing.T, email string) string {
	return fmt.Sprintf(`
data "infra_groups" "%[1]s" {
	filter {
		user_name = "%[2]s"
	}

	include_users = true
}`, t.Name(), email)
}
