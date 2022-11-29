package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/infrahq/infra/uid"
)

func TestAccResourceUserGrant(t *testing.T) {
	var id1, id2, id3, id4 uid.ID

	email := randomEmail()

	cluster := randomName("cluster")
	namespace := randomName("ns")

	resourceName := "infra_user_grant.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserGrant(email, "admin", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "user_email", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "admin"),
				),
			},
			{
				Config: testAccResourceUserGrant(email, "view", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "user_email", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: testAccResourceUserGrant_byName(email, "edit", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "user_email", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "edit"),
					testAccCheckIDChanged(&id2, &id3),
				),
			},
			{
				Config: testAccResourceUserGrant_withNamespace(email, "cluster-admin", cluster, namespace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id4)),
					resource.TestCheckResourceAttr(resourceName, "user_email", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "cluster-admin"),
					testAccCheckIDChanged(&id3, &id4),
				),
			},
		},
	})
}

func testAccResourceUserGrant(email, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	email = "%[1]s"
}

resource "infra_user_grant" "test" {
	user_id = infra_user.test.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
	}
}`, email, role, cluster)
}

func testAccResourceUserGrant_byName(email, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	email = "%[1]s"
}

resource "infra_user_grant" "test" {
	user_email = "%[1]s"

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
	}

	depends_on = [
		infra_user.test,
	]
}`, email, role, cluster)
}

func testAccResourceUserGrant_withNamespace(email, role, cluster, namespace string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	email = "%[1]s"
}

resource "infra_user_grant" "test" {
	user_id = infra_user.test.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
		namespace = "%[4]s"
	}
}`, email, role, cluster, namespace)
}
