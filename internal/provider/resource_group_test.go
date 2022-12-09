package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/infrahq/infra/uid"
)

func TestAccResourceGroup(t *testing.T) {
	var id1, id2, id3 uid.ID

	name1 := randomName()
	name2 := randomName()
	nameWithSpace := fmt.Sprintf("%s %s", randomName(), randomName())

	resourceName := "infra_group.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroup(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "name", name1),
				),
			},
			{
				Config: testAccResourceGroup(name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "name", name2),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: testAccResourceGroup(nameWithSpace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "name", nameWithSpace),
					testAccCheckIDChanged(&id2, &id3),
				),
			},
		},
	})
}

func randomName(prefixes ...string) string {
	prefixes = append([]string{"tf"}, prefixes...)
	return acctest.RandomWithPrefix(strings.Join(prefixes, "-"))
}

func testAccResourceGroup(name string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}`, name)
}
