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

	resourceName := fmt.Sprintf("infra_grant.%s", t.Name())

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGrant_userKubernetes(t, "admin", cluster),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "admin"),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGrant_userKubernetes(t, "view", cluster),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGrant_userKubernetesByName(t, email, "edit", cluster),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "edit"),
					testAccCheckIDChanged(&id2, &id3),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGrant_userKubernetesWithNamespace(t, "cluster-admin", cluster, namespace),
				),
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

func testAccResourceGrant_userKubernetes(t *testing.T, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_grant" "%[1]s" {
	user_id = infra_user.%[1]s.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
	}
}`, t.Name(), role, cluster)
}

func testAccResourceGrant_userKubernetesByName(t *testing.T, email, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_grant" "%[1]s" {
	user_name = "%[2]s"

	kubernetes {
		role = "%[3]s"
		cluster = "%[4]s"
	}

	depends_on = [
		infra_user.%[1]s,
	]
}`, t.Name(), email, role, cluster)
}

func testAccResourceGrant_userKubernetesWithNamespace(t *testing.T, role, cluster, namespace string) string {
	return fmt.Sprintf(`
resource "infra_grant" "%[1]s" {
	user_id = infra_user.%[1]s.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
		namespace = "%[4]s"
	}
}`, t.Name(), role, cluster, namespace)
}

func TestAccResourceGrant_userInfra(t *testing.T) {
	var id1, id2, id3 uid.ID

	email := randomEmail()

	resourceName := fmt.Sprintf("infra_grant.%s", t.Name())

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGrant_userInfra(t, "admin"),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "admin"),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGrant_userInfra(t, "view"),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "user_name", email),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceUser(t, email),
					testAccResourceGrant_userInfraByName(t, email, "admin"),
				),
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

func testAccResourceGrant_userInfra(t *testing.T, role string) string {
	return fmt.Sprintf(`
resource "infra_grant" "%[1]s" {
	user_id = infra_user.%[1]s.id

	infra {
		role = "%[2]s"
	}
}`, t.Name(), role)
}

func testAccResourceGrant_userInfraByName(t *testing.T, email, role string) string {
	return fmt.Sprintf(`
resource "infra_grant" "%[1]s" {
	user_name = "%[2]s"

	infra {
		role = "%[3]s"
	}

	depends_on = [
		infra_user.%[1]s,
	]
}`, t.Name(), email, role)
}

func TestAccResourceGrant_groupKubernetes(t *testing.T) {
	var id1, id2, id3, id4 uid.ID

	email := randomEmail()

	cluster := randomName("cluster")
	namespace := randomName("ns")

	resourceName := fmt.Sprintf("infra_grant.%s", t.Name())

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: composeTestConfigFunc(
					testAccResourceGroup(t, email),
					testAccResourceGrant_groupKubernetes(t, "admin", cluster),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "group_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "admin"),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceGroup(t, email),
					testAccResourceGrant_groupKubernetes(t, "view", cluster),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "group_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceGroup(t, email),
					testAccResourceGrant_groupKubernetesByName(t, email, "edit", cluster),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "group_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "edit"),
					testAccCheckIDChanged(&id2, &id3),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceGroup(t, email),
					testAccResourceGrant_groupKubernetesWithNamespace(t, "cluster-admin", cluster, namespace),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id4)),
					resource.TestCheckResourceAttr(resourceName, "group_name", email),
					resource.TestCheckResourceAttr(resourceName, "kubernetes.0.role", "cluster-admin"),
					testAccCheckIDChanged(&id3, &id4),
				),
			},
		},
	})
}

func testAccResourceGrant_groupKubernetes(t *testing.T, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_grant" "%[1]s" {
	group_id = infra_group.%[1]s.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
	}
}`, t.Name(), role, cluster)
}

func testAccResourceGrant_groupKubernetesByName(t *testing.T, email, role, cluster string) string {
	return fmt.Sprintf(`
resource "infra_grant" "%[1]s" {
	group_name = "%[2]s"

	kubernetes {
		role = "%[3]s"
		cluster = "%[4]s"
	}

	depends_on = [
		infra_group.%[1]s,
	]
}`, t.Name(), email, role, cluster)
}

func testAccResourceGrant_groupKubernetesWithNamespace(t *testing.T, role, cluster, namespace string) string {
	return fmt.Sprintf(`
resource "infra_grant" "%[1]s" {
	group_id = infra_group.%[1]s.id

	kubernetes {
		role = "%[2]s"
		cluster = "%[3]s"
		namespace = "%[4]s"
	}
}`, t.Name(), role, cluster, namespace)
}

func TestAccResourceGrant_groupInfra(t *testing.T) {
	var id1, id2, id3 uid.ID

	email := randomEmail()

	resourceName := fmt.Sprintf("infra_grant.%s", t.Name())

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviders(t),
		Steps: []resource.TestStep{
			{
				Config: composeTestConfigFunc(
					testAccResourceGroup(t, email),
					testAccResourceGrant_groupInfra(t, "admin"),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
					resource.TestCheckResourceAttr(resourceName, "group_name", email),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "admin"),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceGroup(t, email),
					testAccResourceGrant_groupInfra(t, "view"),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
					resource.TestCheckResourceAttr(resourceName, "group_name", email),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "view"),
					testAccCheckIDChanged(&id1, &id2),
				),
			},
			{
				Config: composeTestConfigFunc(
					testAccResourceGroup(t, email),
					testAccResourceGrant_groupInfraByName(t, email, "admin"),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id3)),
					resource.TestCheckResourceAttr(resourceName, "group_name", email),
					resource.TestCheckResourceAttr(resourceName, "infra.0.role", "admin"),
					testAccCheckIDChanged(&id2, &id3),
				),
			},
		},
	})
}

func testAccResourceGrant_groupInfra(t *testing.T, role string) string {
	return fmt.Sprintf(`
resource "infra_grant" "%[1]s" {
	group_id = infra_group.%[1]s.id

	infra {
		role = "%[2]s"
	}
}`, t.Name(), role)
}

func testAccResourceGrant_groupInfraByName(t *testing.T, email, role string) string {
	return fmt.Sprintf(`
resource "infra_grant" "%[1]s" {
	group_name = "%[2]s"

	infra {
		role = "%[3]s"
	}

	depends_on = [
		infra_group.%[1]s,
	]
}`, t.Name(), email, role)
}
