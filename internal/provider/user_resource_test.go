package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	// Test cases for user resource
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tableau_user" "uat_test" {
	email 		 = "uat_test@example.com"
	site_role 	 = "Unlicensed"
	auth_setting = "OpenID"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tableau_user.uat_test", "email", "uat_test@example.com"),
					resource.TestCheckResourceAttr("tableau_user.uat_test", "site_role", "Unlicensed"),
					resource.TestCheckResourceAttr("tableau_user.uat_test", "auth_setting", "OpenID"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("tableau_user.uat_test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "tableau_user.uat_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "tableau_user" "uat_test" {
	email 		 = "uat_test@example.com"
	site_role 	 = "Viewer"
	auth_setting = "SAML"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tableau_user.uat_test", "email", "uat_test@example.com"),
					resource.TestCheckResourceAttr("tableau_user.uat_test", "site_role", "Viewer"),
					resource.TestCheckResourceAttr("tableau_user.uat_test", "auth_setting", "SAML"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
