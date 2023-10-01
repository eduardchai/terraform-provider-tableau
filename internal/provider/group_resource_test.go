package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupResource(t *testing.T) {
	// Test cases for group data source
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tableau_group" "terraform_provider_test" {
	name = "terraform-provider-test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tableau_group.terraform_provider_test", "name", "terraform-provider-test"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("tableau_group.terraform_provider_test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "tableau_group.terraform_provider_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "tableau_group" "terraform_provider_test" {
	name = "terraform-provider-test-updated"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tableau_group.terraform_provider_test", "name", "terraform-provider-test-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
