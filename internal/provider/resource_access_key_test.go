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

	email := randomEmail()

	resourceName := "infra_access_key.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAccessKey(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestMatchResourceAttr(resourceName, "secret", regexp.MustCompile("^[[:alnum:]]{10}\\.[[:alnum:]]{24}$")),
				),
			},
			{
				Config: testAccResourceAccessKey(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestMatchResourceAttr(resourceName, "secret", regexp.MustCompile("^[[:alnum:]]{10}\\.[[:alnum:]]{24}$")),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
		},
	})
}

func testAccResourceAccessKey(email string) string {
	return fmt.Sprintf(`

resource "infra_user" "test" {
	email = "%[1]s"
}

resource "infra_access_key" "test" {
	user_id = infra_user.test.id
}
`, email)
}
