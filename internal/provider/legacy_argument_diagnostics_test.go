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
  hostname = "https://awx.example.invalid"
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

func TestTerraformRejectsUnexpectedRelationshipArguments(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "awx" {
  hostname = "https://awx.example.invalid"
  username = "demo"
  password = "demo"
}

resource "awx_team_user_association" "legacy" {
  team_id          = 1
  user_id          = 2
  legacy_source_id = 1
  legacy_target_id = 2
}
`,
				ExpectError: regexp.MustCompile(`(?s)(Unsupported argument|An argument named "(legacy_source_id|legacy_target_id)" is not expected here)`),
			},
		},
	})
}

func TestTerraformRejectsUnexpectedSurveySpecArguments(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "awx" {
  hostname = "https://awx.example.invalid"
  username = "demo"
  password = "demo"
}

resource "awx_job_template_survey_spec" "legacy" {
  job_template_id = 1
  legacy_id       = 1
  spec = {
    name        = "legacy"
    description = "legacy"
    spec        = []
  }
}
`,
				ExpectError: regexp.MustCompile(`(?s)(Unsupported argument|An argument named "legacy_id" is not expected here)`),
			},
		},
	})
}
