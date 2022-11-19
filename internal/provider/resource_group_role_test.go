package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/infrahq/infra/uid"
)

func TestAccResourceGroupRole(t *testing.T) {
	var id1, id2, id3 uid.ID

	name := randomName()
	resourceName := "infra_group_role.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupRole(name, "admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "role", "admin"),
				),
			},
			{
				Config: testAccResourceGroupRole(name, "view"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: testAccResourceGroupRole_byName(name, "view"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "role", "view"),
				),
			},
		},
	})
}

func testAccResourceGroupRole(name, role string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}

resource "infra_group_role" "test" {
	group_id = infra_group.test.id
	role = "%[2]s"
}`, name, role)
}

func testAccResourceGroupRole_byName(name, role string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}

resource "infra_group_role" "test" {
	group_name = "%[1]s"
	role = "%[2]s"

	depends_on = [
		infra_group.test,
	]
}`, name, role)
}
