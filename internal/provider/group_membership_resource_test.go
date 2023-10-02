package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupMembershipResource(t *testing.T) {
	// Test cases for group data source
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tableau_group" "uat_terraform_provider_test" {
	name = "uat-terraform-provider-test"
}

resource "tableau_user" "uat_test" {
	email 		 = "uat_test@example.com"
	site_role 	 = "Unlicensed"
	auth_setting = "OpenID"
}

resource "tableau_group_membership" "uat_test_group_membership" {
	group_id = tableau_group.uat_terraform_provider_test.id
	users = [
		tableau_user.uat_test.email,
	]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("tableau_group_membership.uat_test_group_membership", "users.*", "uat_test@example.com"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("tableau_group_membership.uat_test_group_membership", "group_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "tableau_group_membership.uat_test_group_membership",
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "group_id",
			},
			// Notes: Update is skipped because it is not possible to mock update at this time.
			//        Errors encountered are:
			//        - Error running post-apply
			//        - Error running post-test destroy

			// Delete testing automatically occurs in TestCase
		},
	})
}
