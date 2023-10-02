package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupDataSource(t *testing.T) {
	// Test cases for user data source
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
resource "tableau_group" "uat_terraform_provider_test" {
	name = "UAT - terraform provider test"
}

data "tableau_group" "uat_terraform_provider_test" {
	name = resource.tableau_group.uat_terraform_provider_test.name
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tableau_group.uat_terraform_provider_test", "name", "UAT - terraform provider test"),
				),
			},
		},
	})
}
