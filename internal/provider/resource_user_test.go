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

	resourceName := fmt.Sprintf("infra_user.%s", t.Name())

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUser(t, email1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "name", email1),
				),
			},
			{
				Config: testAccResourceUser_password(t, email1, "password"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "name", email1),
					resource.TestCheckResourceAttr(resourceName, "password", "password"),
				),
			},
			{
				Config: testAccResourceUser(t, email2),
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

func testAccResourceUser(t *testing.T, email string) string {
	return fmt.Sprintf(`
resource "infra_user" "%[1]s" {
	name = "%[2]s"
}`, t.Name(), email)
}

func testAccResourceUser_password(t *testing.T, email, password string) string {
	return fmt.Sprintf(`
resource "infra_user" "%[1]s" {
	name = "%[2]s"
	password = "%[3]s"
}`, t.Name(), email, password)
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
