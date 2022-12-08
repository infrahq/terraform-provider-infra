package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceGroupMembership(t *testing.T) {
	email := randomEmail()
	name := randomName()

	resourceName := "infra_group_membership.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupMembership(email, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
				),
			},
			{
				Config: testAccResourceGroupMembership_byUserEmail(email, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
				),
			},
			{
				Config: testAccResourceGroupMembership_byGroupName(email, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
				),
			},
		},
	})
}

func testAccResourceGroupMembership(email, name string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	name = "%[1]s"
}

resource "infra_group" "test" {
	name = "%[2]s"
}

resource "infra_group_membership" "test" {
	user_id = infra_user.test.id
	group_id = infra_group.test.id
}`, email, name)
}

func testAccResourceGroupMembership_byUserEmail(email, name string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	name = "%[1]s"
}

resource "infra_group" "test" {
	name = "%[2]s"
}

resource "infra_group_membership" "test" {
	user_name = "%[1]s"
	group_id = infra_group.test.id

	depends_on = [
		infra_user.test,
	]
}`, email, name)
}

func testAccResourceGroupMembership_byGroupName(email, name string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	name = "%[1]s"
}

resource "infra_group" "test" {
	name = "%[2]s"
}

resource "infra_group_membership" "test" {
	user_id = infra_user.test.id
	group_name = "%[2]s"

	depends_on = [
		infra_group.test,
	]
}`, email, name)
}
