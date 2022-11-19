package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/infrahq/infra/uid"
)

func TestAccResourceGroupGrant(t *testing.T) {
	var id1, id2, id3, id4 uid.ID

	name := randomName()

	cluster := randomName("cluster")
	namespace := randomName("ns")

	resourceName := "infra_group_grant.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupGrant(name, "admin", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "admin"),
				),
			},
			{
				Config: testAccResourceGroupGrant(name, "view", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: testAccResourceGroupGrant_byName(name, "edit", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "edit"),
					testAccCheckIDChanged(&id2, &id3),
				),
			},
			{
				Config: testAccResourceGroupGrant_withNamespace(name, "cluster-admin", cluster, namespace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id4)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "cluster-admin"),
					testAccCheckIDChanged(&id3, &id4),
				),
			},
		},
	})
}

func testAccResourceGroupGrant(name, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}

resource "infra_group_grant" "test" {
	group_id = infra_group.test.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
	}
}`, name, role, cluster)
}

func testAccResourceGroupGrant_byName(name, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}

resource "infra_group_grant" "test" {
	group_name = "%[1]s"

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
	}

	depends_on = [
		infra_group.test,
	]
}`, name, role, cluster)
}

func testAccResourceGroupGrant_withNamespace(name, role, cluster, namespace string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}

resource "infra_group_grant" "test" {
	group_id = infra_group.test.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
		namespace = "%[4]s"
	}
}`, name, role, cluster, namespace)
}
