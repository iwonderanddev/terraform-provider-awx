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
	t.Logf("starting terraform acceptance: team CRUD/import team=%q organization=%d", teamName, organizationID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreStep(t, "step 1/4: apply create config without description for team=%q", teamName),
				Config:    configCreateNoDescription,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLog(t, "step 1/4 complete: team created without description"),
					resource.TestCheckResourceAttr(resourceName, "name", teamName),
					resource.TestCheckResourceAttr(resourceName, "organization", strconv.FormatInt(organizationID, 10)),
					resource.TestCheckNoResourceAttr(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				PreConfig: testAccPreStep(t, "step 2/4: apply update to set description for team=%q", teamName),
				Config:    configCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLog(t, "step 2/4 complete: description set"),
					resource.TestCheckResourceAttr(resourceName, "name", teamName),
					resource.TestCheckResourceAttr(resourceName, "organization", strconv.FormatInt(organizationID, 10)),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform-plugin-testing"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				PreConfig: testAccPreStep(t, "step 3/4: apply update to change description for team=%q", teamName),
				Config:    configUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLog(t, "step 3/4 complete: description changed"),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform-plugin-testing"),
				),
			},
			{
				PreConfig:         testAccPreStep(t, "step 4/4: import team resource state"),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcceptanceTerraform_OrganizationResourceDefaultedFieldStability(t *testing.T) {
	t.Parallel()

	_ = testAccPreCheck(t, envOrganizationID)
	organizationName := fmt.Sprintf("tf-awx-org-%d", time.Now().UnixNano())
	resourceName := "awx_organization.test"
	config := testAccOrganizationResourceConfigWithoutMaxHosts(organizationName)
	t.Logf("starting terraform acceptance: organization default-field stability organization=%q", organizationName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreStep(t, "step 1/3: apply create config without max_hosts for organization=%q", organizationName),
				Config:    config,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLog(t, "step 1/3 complete: organization created"),
					resource.TestCheckResourceAttr(resourceName, "name", organizationName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				PreConfig: testAccPreStep(t, "step 2/3: run plan-only to assert no drift for omitted max_hosts"),
				Config:    config,
				PlanOnly:  true,
			},
			{
				PreConfig:         testAccPreStep(t, "step 3/3: import organization resource state"),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcceptanceTerraform_InventoryResourceDefaultedFieldStability(t *testing.T) {
	t.Parallel()

	_ = testAccPreCheck(t, envOrganizationID)
	organizationName := fmt.Sprintf("tf-awx-inv-org-%d", time.Now().UnixNano())
	inventoryName := fmt.Sprintf("tf-awx-inv-%d", time.Now().UnixNano())
	resourceName := "awx_inventory.test"
	config := testAccInventoryResourceConfigWithoutFallback(inventoryName, organizationName)
	t.Logf("starting terraform acceptance: inventory default-field stability inventory=%q organization=%q", inventoryName, organizationName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreStep(t, "step 1/3: apply create config without prevent_instance_group_fallback for inventory=%q", inventoryName),
				Config:    config,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLog(t, "step 1/3 complete: inventory created and default fallback value captured"),
					resource.TestCheckResourceAttr(resourceName, "name", inventoryName),
					resource.TestCheckResourceAttr(resourceName, "prevent_instance_group_fallback", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				PreConfig: testAccPreStep(t, "step 2/3: run plan-only to assert no drift for omitted prevent_instance_group_fallback"),
				Config:    config,
				PlanOnly:  true,
			},
			{
				PreConfig:         testAccPreStep(t, "step 3/3: import inventory resource state"),
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
	t.Logf("starting terraform acceptance: team data source lookup team=%q organization=%d", teamName, organizationID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreStep(t, "step 1/1: apply team resource and lookup by name"),
				Config:    testAccTeamDataSourceConfig(teamName, organizationID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLog(t, "step 1/1 complete: team data source resolved id and name"),
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
	t.Logf("starting terraform acceptance: team-user relationship team=%q organization=%d user=%d", teamName, organizationID, userID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreStep(t, "step 1/2: apply team and relationship association for user=%d", userID),
				Config:    testAccRelationshipConfig(teamName, organizationID, userID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLog(t, "step 1/2 complete: relationship association exists"),
					resource.TestCheckResourceAttr(resourceName, "child_id", strconv.FormatInt(userID, 10)),
					testCheckCompositeRelationshipID(resourceName, teamResourceName, userID),
				),
			},
			{
				PreConfig:    testAccPreStep(t, "step 2/2: import relationship using composite id"),
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

func testAccPreStep(t *testing.T, format string, args ...any) func() {
	t.Helper()
	return func() {
		t.Helper()
		t.Logf(format, args...)
	}
}

func testAccCheckLog(t *testing.T, format string, args ...any) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		t.Helper()
		t.Logf(format, args...)
		return nil
	}
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

func testAccOrganizationResourceConfigWithoutMaxHosts(name string) string {
	return fmt.Sprintf(`
%s
resource "awx_organization" "test" {
  name        = %q
  description = "created for default field stability checks"
}
`, testAccProviderConfig(), name)
}

func testAccInventoryResourceConfigWithoutFallback(inventoryName string, organizationName string) string {
	return fmt.Sprintf(`
%s
resource "awx_organization" "test" {
  name        = %q
  description = "created for inventory default field stability checks"
}

resource "awx_inventory" "test" {
  name         = %q
  organization = tonumber(awx_organization.test.id)
  description  = "created for inventory default field stability checks"
}
`, testAccProviderConfig(), organizationName, inventoryName)
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
