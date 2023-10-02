package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	// Test cases for user data source
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
resource "tableau_user" "uat_test" {
	email 		 = "uat_test@example.com"
	site_role 	 = "Unlicensed"
	auth_setting = "OpenID"
}

data "tableau_user" "uat_test" {
	email = resource.tableau_user.uat_test.email
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tableau_user.uat_test", "email", "uat_test@example.com"),
					resource.TestCheckResourceAttr("data.tableau_user.uat_test", "site_role", "Unlicensed"),
					resource.TestCheckResourceAttr("data.tableau_user.uat_test", "auth_setting", "OpenID"),
				),
			},
		},
	})
}
