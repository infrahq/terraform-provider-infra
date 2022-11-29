package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceUserGroup(t *testing.T) {
	// var id1, id2 uid.ID

	email := randomEmail()
	name := randomName()

	resourceName := "infra_user_group.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserGroup(email, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_email", email),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
				),
			},
			{
				Config: testAccResourceUserGroup_byUserEmail(email, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_email", email),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
				),
			},
			{
				Config: testAccResourceUserGroup_byGroupName(email, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_email", email),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
				),
			},
		},
	})
}

func testAccResourceUserGroup(email, name string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	email = "%[1]s"
}

resource "infra_group" "test" {
	name = "%[2]s"
}

resource "infra_user_group" "test" {
	user_id = infra_user.test.id
	group_id = infra_group.test.id
}`, email, name)
}

func testAccResourceUserGroup_byUserEmail(email, name string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	email = "%[1]s"
}

resource "infra_group" "test" {
	name = "%[2]s"
}

resource "infra_user_group" "test" {
	user_email = "%[1]s"
	group_id = infra_group.test.id

	depends_on = [
		infra_user.test,
	]
}`, email, name)
}

func testAccResourceUserGroup_byGroupName(email, name string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	email = "%[1]s"
}

resource "infra_group" "test" {
	name = "%[2]s"
}

resource "infra_user_group" "test" {
	user_id = infra_user.test.id
	group_name = "%[2]s"

	depends_on = [
		infra_group.test,
	]
}`, email, name)
}
