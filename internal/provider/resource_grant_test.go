package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/infrahq/infra/uid"
)

func TestAccResourceGrant_userKubernetes(t *testing.T) {
	var id1, id2, id3, id4 uid.ID

	email := randomEmail()

	cluster := randomName("cluster")
	namespace := randomName("ns")

	resourceName := "infra_grant.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGrant_userKubernetes(email, "admin", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "admin"),
				),
			},
			{
				Config: testAccResourceGrant_userKubernetes(email, "view", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: testAccResourceGrant_userKubernetesByName(email, "edit", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "edit"),
					testAccCheckIDChanged(&id2, &id3),
				),
			},
			{
				Config: testAccResourceGrant_userKubernetesWithNamespace(email, "cluster-admin", cluster, namespace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id4)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "cluster-admin"),
					testAccCheckIDChanged(&id3, &id4),
				),
			},
		},
	})
}

func testAccResourceGrant_userKubernetes(email, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	name = "%[1]s"
}

resource "infra_grant" "test" {
	user_id = infra_user.test.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
	}
}`, email, role, cluster)
}

func testAccResourceGrant_userKubernetesByName(email, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	name = "%[1]s"
}

resource "infra_grant" "test" {
	user_name = "%[1]s"

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
	}

	depends_on = [
		infra_user.test,
	]
}`, email, role, cluster)
}

func testAccResourceGrant_userKubernetesWithNamespace(email, role, cluster, namespace string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	name = "%[1]s"
}

resource "infra_grant" "test" {
	user_id = infra_user.test.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
		namespace = "%[4]s"
	}
}`, email, role, cluster, namespace)
}

func TestAccResourceGrant_userInfra(t *testing.T) {
	var id1, id2, id3 uid.ID

	email := randomEmail()

	cluster := randomName("cluster")

	resourceName := "infra_grant.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGrant_userInfra(email, "admin", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "admin"),
				),
			},
			{
				Config: testAccResourceGrant_userInfra(email, "view", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: testAccResourceGrant_userInfraByName(email, "admin", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "admin"),
					testAccCheckIDChanged(&id2, &id3),
				),
			},
		},
	})
}

func testAccResourceGrant_userInfra(email, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	name = "%[1]s"
}

resource "infra_grant" "test" {
	user_id = infra_user.test.id

	infra {
		role = "%[2]s"
	}
}`, email, role, cluster)
}

func testAccResourceGrant_userInfraByName(email, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_user" "test" {
	name = "%[1]s"
}

resource "infra_grant" "test" {
	user_name = "%[1]s"

	infra {
		role = "%[2]s"
	}

	depends_on = [
		infra_user.test,
	]
}`, email, role, cluster)
}

func TestAccResourceGrant_groupKubernetes(t *testing.T) {
	var id1, id2, id3, id4 uid.ID

	name := randomName()

	cluster := randomName("cluster")
	namespace := randomName("ns")

	resourceName := "infra_grant.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGrant_groupKubernetes(name, "admin", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "admin"),
				),
			},
			{
				Config: testAccResourceGrant_groupKubernetes(name, "view", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: testAccResourceGrant_groupKubernetesByName(name, "edit", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "edit"),
					testAccCheckIDChanged(&id2, &id3),
				),
			},
			{
				Config: testAccResourceGrant_groupKubernetesWithNamespace(name, "cluster-admin", cluster, namespace),
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

func testAccResourceGrant_groupKubernetes(name, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}

resource "infra_grant" "test" {
	group_id = infra_group.test.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
	}
}`, name, role, cluster)
}

func testAccResourceGrant_groupKubernetesByName(name, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}

resource "infra_grant" "test" {
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

func testAccResourceGrant_groupKubernetesWithNamespace(name, role, cluster, namespace string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}

resource "infra_grant" "test" {
	group_id = infra_group.test.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
		namespace = "%[4]s"
	}
}`, name, role, cluster, namespace)
}

func TestAccResourceGrant_groupInfra(t *testing.T) {
	var id1, id2, id3 uid.ID

	name := randomName()

	cluster := randomName("cluster")

	resourceName := "infra_grant.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGrant_groupInfra(name, "admin", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "admin"),
				),
			},
			{
				Config: testAccResourceGrant_groupInfra(name, "view", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: testAccResourceGrant_groupInfraByName(name, "admin", cluster),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "group_name", name),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "admin"),
					testAccCheckIDChanged(&id2, &id3),
				),
			},
		},
	})
}

func testAccResourceGrant_groupInfra(name, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}

resource "infra_grant" "test" {
	group_id = infra_group.test.id

	infra {
		role = "%[2]s"
	}
}`, name, role, cluster)
}

func testAccResourceGrant_groupInfraByName(name, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_group" "test" {
	name = "%[1]s"
}

resource "infra_grant" "test" {
	group_name = "%[1]s"

	infra {
		role = "%[2]s"
	}

	depends_on = [
		infra_group.test,
	]
}`, name, role, cluster)
}
