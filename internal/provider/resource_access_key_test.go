package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/infrahq/infra/uid"
)

func TestAccResourceAccessKey(t *testing.T) {
	var id1, id2 uid.ID

	resourceName := fmt.Sprintf("infra_access_key.%s", t.Name())

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAccessKey_connector(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestMatchResourceAttr(resourceName, "secret", regexp.MustCompile("^[[:alnum:]]{10}\\.[[:alnum:]]{24}$")),
				),
			},
			{
				Config: testAccResourceAccessKey_connector(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestMatchResourceAttr(resourceName, "secret", regexp.MustCompile("^[[:alnum:]]{10}\\.[[:alnum:]]{24}$")),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
		},
	})
}

func testAccResourceAccessKey_connector(t *testing.T) string {
	return fmt.Sprintf(`
resource "infra_access_key" "%[1]s" {}
`, t.Name())
}
