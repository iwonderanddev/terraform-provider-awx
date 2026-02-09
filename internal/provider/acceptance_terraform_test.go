package provider_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	awxprovider "github.com/damien/terraform-awx-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	envAcceptance      = "AWX_ACCEPTANCE"
	envBaseURL         = "AWX_BASE_URL"
	envUsername        = "AWX_USERNAME"
	envPassword        = "AWX_PASSWORD"
	envOrganizationID  = "AWX_TEST_ORGANIZATION_ID"
	envRelationshipUID = "AWX_TEST_USER_ID"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"awx": providerserver.NewProtocol6WithError(awxprovider.New("test")()),
}

func TestAcceptanceTerraform_TeamResourceCRUDAndImport(t *testing.T) {
	t.Parallel()

	organizationID := testAccPreCheck(t, envOrganizationID)
	teamName := fmt.Sprintf("tf-awx-team-%d", time.Now().UnixNano())

	resourceName := "awx_team.test"
	configCreateNoDescription := testAccTeamResourceConfigWithoutDescription(teamName, organizationID)
	configCreate := testAccTeamResourceConfig(teamName, organizationID, "created by terraform-plugin-testing")
	configUpdate := testAccTeamResourceConfig(teamName, organizationID, "updated by terraform-plugin-testing")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configCreateNoDescription,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", teamName),
					resource.TestCheckResourceAttr(resourceName, "organization", strconv.FormatInt(organizationID, 10)),
					resource.TestCheckNoResourceAttr(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: configCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", teamName),
					resource.TestCheckResourceAttr(resourceName, "organization", strconv.FormatInt(organizationID, 10)),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform-plugin-testing"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform-plugin-testing"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcceptanceTerraform_TeamDataSourceLookup(t *testing.T) {
	t.Parallel()

	organizationID := testAccPreCheck(t, envOrganizationID)
	teamName := fmt.Sprintf("tf-awx-ds-team-%d", time.Now().UnixNano())

	resourceName := "awx_team.test"
	dataSourceName := "data.awx_team.by_name"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamDataSourceConfig(teamName, organizationID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", teamName),
					resource.TestCheckResourceAttr(dataSourceName, "name", teamName),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
				),
			},
		},
	})
}

func TestAcceptanceTerraform_TeamUserRelationshipResource(t *testing.T) {
	t.Parallel()

	organizationID := testAccPreCheck(t, envOrganizationID)
	userIDRaw := os.Getenv(envRelationshipUID)
	if userIDRaw == "" {
		t.Skipf("Missing %s for relationship acceptance scenario", envRelationshipUID)
	}
	userID, err := strconv.ParseInt(userIDRaw, 10, 64)
	if err != nil {
		t.Fatalf("invalid %s: %v", envRelationshipUID, err)
	}

	teamName := fmt.Sprintf("tf-awx-rel-team-%d", time.Now().UnixNano())
	resourceName := "awx_team_user_association.membership"
	teamResourceName := "awx_team.parent"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRelationshipConfig(teamName, organizationID, userID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "child_id", strconv.FormatInt(userID, 10)),
					testCheckCompositeRelationshipID(resourceName, teamResourceName, userID),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					teamResource, ok := state.RootModule().Resources[teamResourceName]
					if !ok || teamResource == nil {
						return "", fmt.Errorf("missing %s in state", teamResourceName)
					}
					return fmt.Sprintf("%s:%d", teamResource.Primary.ID, userID), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPreCheck(t *testing.T, required ...string) int64 {
	t.Helper()

	if os.Getenv(envAcceptance) != "1" {
		t.Skipf("Acceptance tests are opt-in. Set %s=1 to run Terraform acceptance tests.", envAcceptance)
	}

	mandatory := []string{envBaseURL, envUsername, envPassword}
	mandatory = append(mandatory, required...)

	for _, key := range mandatory {
		if os.Getenv(key) == "" {
			t.Skipf("Missing required acceptance environment variable: %s", key)
		}
	}

	organizationIDRaw := os.Getenv(envOrganizationID)
	organizationID, err := strconv.ParseInt(organizationIDRaw, 10, 64)
	if err != nil {
		t.Fatalf("invalid %s: %v", envOrganizationID, err)
	}

	return organizationID
}

func testAccProviderConfig() string {
	return fmt.Sprintf(`provider "awx" {
  base_url = %q
  username = %q
  password = %q
}
`, os.Getenv(envBaseURL), os.Getenv(envUsername), os.Getenv(envPassword))
}

func testAccTeamResourceConfigWithoutDescription(name string, organizationID int64) string {
	return fmt.Sprintf(`
%s
resource "awx_team" "test" {
  name         = %q
  organization = %d
}
`, testAccProviderConfig(), name, organizationID)
}

func testAccTeamResourceConfig(name string, organizationID int64, description string) string {
	return fmt.Sprintf(`
%s
resource "awx_team" "test" {
  name         = %q
  organization = %d
  description  = %q
}
`, testAccProviderConfig(), name, organizationID, description)
}

func testAccTeamDataSourceConfig(name string, organizationID int64) string {
	return fmt.Sprintf(`
%s
resource "awx_team" "test" {
  name         = %q
  organization = %d
  description  = "created for data source lookup"
}

data "awx_team" "by_name" {
  name = awx_team.test.name
}
`, testAccProviderConfig(), name, organizationID)
}

func testAccRelationshipConfig(name string, organizationID int64, userID int64) string {
	return fmt.Sprintf(`
%s
resource "awx_team" "parent" {
  name         = %q
  organization = %d
  description  = "created for relationship acceptance"
}

resource "awx_team_user_association" "membership" {
  parent_id = tonumber(awx_team.parent.id)
  child_id  = %d
}
`, testAccProviderConfig(), name, organizationID, userID)
}

func testCheckCompositeRelationshipID(relationshipResourceName string, teamResourceName string, userID int64) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		relationshipResource, ok := state.RootModule().Resources[relationshipResourceName]
		if !ok || relationshipResource == nil {
			return fmt.Errorf("missing %s in state", relationshipResourceName)
		}

		teamResource, ok := state.RootModule().Resources[teamResourceName]
		if !ok || teamResource == nil {
			return fmt.Errorf("missing %s in state", teamResourceName)
		}

		expectedID := fmt.Sprintf("%s:%d", teamResource.Primary.ID, userID)
		if relationshipResource.Primary.ID != expectedID {
			return fmt.Errorf("unexpected relationship composite ID: got=%s want=%s", relationshipResource.Primary.ID, expectedID)
		}

		return nil
	}
}
