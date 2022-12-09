package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceGroupMembership(t *testing.T) {
	email := randomEmail()
	name := randomName()

	resourceName := fmt.Sprintf("infra_group_membership.%s", t.Name())

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
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGroup(t, name),
					testAccResourceGroupMembership_byUserEmail(t, email),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGroup(t, name),
					testAccResourceGroupMembership_byGroupName(t, name),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
				),
			},
		},
	})
}

func testAccResourceGroupMembership(t *testing.T) string {
	return fmt.Sprintf(`
resource "infra_group_membership" "%[1]s" {
	user_id = infra_user.%[1]s.id
	group_id = infra_group.%[1]s.id
}`, t.Name())
}

func testAccResourceGroupMembership_byUserEmail(t *testing.T, email string) string {
	return fmt.Sprintf(`
resource "infra_group_membership" "%[1]s" {
	user_name = "%[2]s"
	group_id = infra_group.%[1]s.id

	depends_on = [
		infra_user.%[1]s,
	]
}`, t.Name(), email)
}

func testAccResourceGroupMembership_byGroupName(t *testing.T, name string) string {
	return fmt.Sprintf(`
resource "infra_group_membership" "%[1]s" {
	user_id = infra_user.%[1]s.id
	group_name = "%[2]s"

	depends_on = [
		infra_group.%[1]s,
	]
}`, t.Name(), name)
}
