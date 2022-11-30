package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/infrahq/infra/uid"
)

func TestAccResourceUser(t *testing.T) {
	var id1, id2 uid.ID

	email1 := randomEmail()
	email2 := randomEmail()

	resourceName := "infra_user.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUser(email1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "name", email1),
				),
			},
			{
				Config: testAccResourceUser_password(email1, "password"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "name", email1),
					resource.TestCheckResourceAttr(resourceName, "password", "password"),
				),
			},
			{
				Config: testAccResourceUser(email2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "name", email2),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
		},
	})
}

func randomEmail() string {
	return fmt.Sprintf("%s@example.com", randomName())
}

func testAccResourceUser(email string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	name = "%[1]s"
}`, email)
}

func testAccResourceUser_password(email, password string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	name = "%[1]s"
	password = "%[2]s"
}`, email, password)
}

func testCheckResourceAttrWithID(out *uid.ID) func(s string) error {
	return func(s string) error {
		id, err := uid.Parse([]byte(s))
		if err != nil {
			return err
		}

		*out = id
		return nil
	}
}

func testAccCheckIDChanged(id1, id2 *uid.ID) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *id1 != *id2 {
			return nil
		}

		return fmt.Errorf("resource should have been recreated")
	}
}

func testAccCheckIDUnchanged(id1, id2 *uid.ID) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *id1 == *id2 {
			return nil
		}

		return fmt.Errorf("resource should not have been recreated")
	}
}
