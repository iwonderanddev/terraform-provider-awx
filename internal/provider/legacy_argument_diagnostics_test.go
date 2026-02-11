package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTerraformRejectsLegacyUnsuffixedReferenceArgument(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "awx" {
  base_url = "https://awx.example.invalid"
  username = "demo"
  password = "demo"
}

resource "awx_team" "legacy" {
  name            = "legacy"
  organization_id = 1
  organization    = 1
}
`,
				ExpectError: regexp.MustCompile(`(?s)(Unsupported argument|An argument named "organization" is not expected here)`),
			},
		},
	})
}
