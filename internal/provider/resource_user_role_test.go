package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/infrahq/infra/uid"
)

func TestAccResourceUserRole(t *testing.T) {
	var id1, id2, id3 uid.ID

	email := randomEmail()
	resourceName := "infra_user_role.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserRole(email, "admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "user_email", email),
					resource.TestCheckResourceAttr(resourceName, "role", "admin"),
				),
			},
			{
				Config: testAccResourceUserRole(email, "view"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "user_email", email),
					resource.TestCheckResourceAttr(resourceName, "role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: testAccResourceUserRole_byEmail(email, "view"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "user_email", email),
					resource.TestCheckResourceAttr(resourceName, "role", "view"),
				),
			},
		},
	})
}

func testAccResourceUserRole(email, role string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	email = "%[1]s"
}

resource "infra_user_role" "test" {
	user_id = infra_user.test.id
	role = "%[2]s"
}`, email, role)
}

func testAccResourceUserRole_byEmail(email, role string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	email = "%[1]s"
}

resource "infra_user_role" "test" {
	user_email = "%[1]s"
	role = "%[2]s"

	depends_on = [
		infra_user.test,
	]
}`, email, role)
}
