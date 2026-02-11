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

func TestTerraformRejectsLegacyRelationshipDirectionalArguments(t *testing.T) {
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

resource "awx_team_user_association" "legacy" {
  team_id   = 1
  user_id   = 2
  parent_id = 1
  child_id  = 2
}
`,
				ExpectError: regexp.MustCompile(`(?s)(Unsupported argument|An argument named "(parent_id|child_id)" is not expected here)`),
			},
		},
	})
}

func TestTerraformRejectsLegacySurveySpecParentArgument(t *testing.T) {
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

resource "awx_job_template_survey_spec" "legacy" {
  job_template_id = 1
  parent_id       = 1
  spec = jsonencode({
    name        = "legacy"
    description = "legacy"
    spec        = []
  })
}
`,
				ExpectError: regexp.MustCompile(`(?s)(Unsupported argument|An argument named "parent_id" is not expected here)`),
			},
		},
	})
}
