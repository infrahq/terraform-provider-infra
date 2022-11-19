package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/infrahq/infra/uid"
)

func TestAccResourceIdentityProvider(t *testing.T) {
	resourceName := "infra_identity_provider.test"

	type testCase struct {
		Name           string
		ConfigFunc     func(string, string, string) string
		ExpectedIssuer string
	}

	testCases := []testCase{
		{
			Name:           "oidc",
			ConfigFunc:     testAccResourceIdentityProvider_withIssuer,
			ExpectedIssuer: "https://my.custom.example.com",
		},
		{
			Name:           "azure",
			ConfigFunc:     testAccResourceIdentityProvider_withAzureAD,
			ExpectedIssuer: "https://login.microsoftonline.com/abc/v2.0",
		},
		{
			Name:           "google",
			ConfigFunc:     testAccResourceIdentityProvider_withGoogle,
			ExpectedIssuer: "https://accounts.google.com",
		},
		{
			Name:           "google-group",
			ConfigFunc:     testAccResourceIdentityProvider_withGoogleGroup,
			ExpectedIssuer: "https://accounts.google.com",
		},
		{
			Name:           "okta",
			ConfigFunc:     testAccResourceIdentityProvider_withOkta,
			ExpectedIssuer: "https://my.okta.example.com",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var id1, id2 uid.ID

			name1 := randomName(testCase.Name)
			name2 := randomName("updated")

			resource.UnitTest(t, resource.TestCase{
				PreCheck:          testAccPreCheck(t),
				ProviderFactories: testAccProviders(t),
				Steps: []resource.TestStep{
					{
						Config: testCase.ConfigFunc(name1, "client_id", "client_secret"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id1)),
							resource.TestCheckResourceAttr(resourceName, "issuer", testCase.ExpectedIssuer),
						),
					},
					{
						Config: testCase.ConfigFunc(name2, "client_id", "client_secret"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
							resource.TestCheckResourceAttr(resourceName, "issuer", testCase.ExpectedIssuer),
							testAccCheckIDUnchanged(&id1, &id2),
						),
					},
					{
						Config: testCase.ConfigFunc(name2, "different_client_id", "client_secret"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
							resource.TestCheckResourceAttr(resourceName, "issuer", testCase.ExpectedIssuer),
							testAccCheckIDUnchanged(&id1, &id2),
						),
					},
					{
						Config: testCase.ConfigFunc(name2, "different_client_id", "different_client_secret"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrWith(resourceName, "id", testCheckResourceAttrWithID(&id2)),
							resource.TestCheckResourceAttr(resourceName, "issuer", testCase.ExpectedIssuer),
							testAccCheckIDUnchanged(&id1, &id2),
						),
					},
				},
			})
		})
	}
}

func testAccResourceIdentityProvider_withIssuer(name, clientID, clientSecret string) string {
	return fmt.Sprintf(`
resource "infra_identity_provider" "test" {
	name = "%[1]s"
	client_id = "%[2]s"
	client_secret = "%[3]s"
	issuer = "https://my.custom.example.com"
}`, name, clientID, clientSecret)
}

func testAccResourceIdentityProvider_withAzureAD(name, clientID, clientSecret string) string {
	return fmt.Sprintf(`
resource "infra_identity_provider" "test" {
	name = "%[1]s"
	client_id = "%[2]s"
	client_secret = "%[3]s"
	azure {
		tenant_id = "abc"
	}
}`, name, clientID, clientSecret)
}

func testAccResourceIdentityProvider_withGoogle(name, clientID, clientSecret string) string {
	return fmt.Sprintf(`
resource "infra_identity_provider" "test" {
	name = "%[1]s"
	client_id = "%[2]s"
	client_secret = "%[3]s"
	google {}
}`, name, clientID, clientSecret)
}

func testAccResourceIdentityProvider_withGoogleGroup(name, clientID, clientSecret string) string {
	return fmt.Sprintf(`
resource "infra_identity_provider" "test" {
	name = "%[1]s"
	client_id = "%[2]s"
	client_secret = "%[3]s"
	google {
		admin_email = "admin@example.com"
		service_account_key = jsonencode({
			"client_email": "client@example.com",
			"private_key": "...",
		})
	}
}`, name, clientID, clientSecret)
}

func testAccResourceIdentityProvider_withOkta(name, clientID, clientSecret string) string {
	return fmt.Sprintf(`
resource "infra_identity_provider" "test" {
	name = "%[1]s"
	client_id = "%[2]s"
	client_secret = "%[3]s"
	okta {
		issuer = "https://my.okta.example.com"
	}
}`, name, clientID, clientSecret)
}
