package acceptance

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/damien/terraform-provider-awx-iwd/internal/client"
)

const (
	envAcceptance       = "AWX_ACCEPTANCE"
	envBaseURL          = "AWX_BASE_URL"
	envUsername         = "AWX_USERNAME"
	envPassword         = "AWX_PASSWORD"
	envOrganizationID   = "AWX_TEST_ORGANIZATION_ID"
	envRelationshipTeam = "AWX_TEST_TEAM_ID"
	envRelationshipUser = "AWX_TEST_USER_ID"
)

var compositeIDPattern = regexp.MustCompile(`^[0-9]+:[0-9]+$`)

func TestAcceptance_TeamCRUDAndImport(t *testing.T) {
	t.Parallel()

	requireAcceptanceEnabled(t)
	env := requireEnv(t, envBaseURL, envUsername, envPassword, envOrganizationID)

	organizationID, err := strconv.ParseInt(env[envOrganizationID], 10, 64)
	if err != nil {
		t.Fatalf("invalid %s: %v", envOrganizationID, err)
	}

	awxClient := mustClient(t, env)
	ctx := context.Background()
	resourceName := fmt.Sprintf("tf-awx-acceptance-team-%d", time.Now().UnixNano())

	t.Logf("creating team %q in organization %d", resourceName, organizationID)
	created, err := awxClient.CreateObject(ctx, "/api/v2/teams/", map[string]any{
		"name":         resourceName,
		"description":  "terraform provider acceptance test",
		"organization": organizationID,
	})
	if err != nil {
		t.Fatalf("failed to create team: %v", err)
	}

	teamID, err := parseID(created["id"])
	if err != nil {
		t.Fatalf("failed to parse created team id: %v", err)
	}
	defer func() {
		t.Logf("cleanup: deleting team id=%d", teamID)
		_ = awxClient.DeleteObject(ctx, "/api/v2/teams/{id}/", strconv.FormatInt(teamID, 10))
	}()

	t.Logf("updating team id=%d description", teamID)
	if _, err := awxClient.UpdateObject(ctx, "/api/v2/teams/{id}/", strconv.FormatInt(teamID, 10), map[string]any{"description": "updated by acceptance test"}); err != nil {
		t.Fatalf("failed to update team: %v", err)
	}

	t.Logf("reading team id=%d", teamID)
	fetched, err := awxClient.GetObject(ctx, "/api/v2/teams/{id}/", strconv.FormatInt(teamID, 10))
	if err != nil {
		t.Fatalf("failed to read team: %v", err)
	}
	if got := fmt.Sprintf("%v", fetched["name"]); got != resourceName {
		t.Fatalf("unexpected team name: got=%q want=%q", got, resourceName)
	}

	importID := strconv.FormatInt(teamID, 10)
	t.Logf("validating import id format: %s", importID)
	if _, err := strconv.ParseInt(importID, 10, 64); err != nil {
		t.Fatalf("expected numeric import id, got %q", importID)
	}

	t.Logf("deleting team id=%d", teamID)
	if err := awxClient.DeleteObject(ctx, "/api/v2/teams/{id}/", strconv.FormatInt(teamID, 10)); err != nil {
		t.Fatalf("failed to delete team: %v", err)
	}

	t.Logf("confirming team id=%d no longer exists", teamID)
	_, err = awxClient.GetObject(ctx, "/api/v2/teams/{id}/", strconv.FormatInt(teamID, 10))
	if err == nil {
		t.Fatalf("expected team lookup to fail after delete")
	}
	apiErr, ok := err.(*client.APIError)
	if !ok || apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 APIError after delete, got %T: %v", err, err)
	}
}

func TestAcceptance_TeamUserRelationshipCompositeImport(t *testing.T) {
	t.Parallel()

	requireAcceptanceEnabled(t)
	env := requireEnv(t, envBaseURL, envUsername, envPassword, envRelationshipTeam, envRelationshipUser)

	teamID, err := strconv.ParseInt(env[envRelationshipTeam], 10, 64)
	if err != nil {
		t.Fatalf("invalid %s: %v", envRelationshipTeam, err)
	}
	userID, err := strconv.ParseInt(env[envRelationshipUser], 10, 64)
	if err != nil {
		t.Fatalf("invalid %s: %v", envRelationshipUser, err)
	}

	awxClient := mustClient(t, env)
	ctx := context.Background()
	relationshipPath := "/api/v2/teams/{id}/users/"

	t.Logf("checking initial team-user relationship state team=%d user=%d", teamID, userID)
	existedBefore, err := awxClient.RelationshipExists(ctx, relationshipPath, teamID, userID)
	if err != nil {
		t.Fatalf("failed pre-check for relationship existence: %v", err)
	}

	t.Logf("associating team=%d with user=%d", teamID, userID)
	if err := awxClient.Associate(ctx, relationshipPath, teamID, userID); err != nil {
		t.Fatalf("failed to associate team and user: %v", err)
	}

	t.Logf("verifying association exists team=%d user=%d", teamID, userID)
	existsAfterCreate, err := awxClient.RelationshipExists(ctx, relationshipPath, teamID, userID)
	if err != nil {
		t.Fatalf("failed relationship read after create: %v", err)
	}
	if !existsAfterCreate {
		t.Fatalf("expected relationship to exist after associate")
	}

	compositeID := fmt.Sprintf("%d:%d", teamID, userID)
	t.Logf("validating composite import id format: %s", compositeID)
	if !compositeIDPattern.MatchString(compositeID) {
		t.Fatalf("invalid composite import id: %q", compositeID)
	}

	t.Logf("disassociating team=%d from user=%d", teamID, userID)
	if err := awxClient.Disassociate(ctx, relationshipPath, teamID, userID); err != nil {
		t.Fatalf("failed to disassociate team and user: %v", err)
	}

	t.Logf("verifying association is removed team=%d user=%d", teamID, userID)
	existsAfterDelete, err := awxClient.RelationshipExists(ctx, relationshipPath, teamID, userID)
	if err != nil {
		t.Fatalf("failed relationship read after delete: %v", err)
	}
	if existsAfterDelete {
		t.Fatalf("expected relationship to be absent after disassociate")
	}

	if existedBefore {
		t.Logf("restoring initial association team=%d user=%d", teamID, userID)
		if err := awxClient.Associate(ctx, relationshipPath, teamID, userID); err != nil {
			t.Fatalf("failed to restore original relationship state: %v", err)
		}
	}
}

func TestAcceptance_TeamLookupByIDAndName(t *testing.T) {
	t.Parallel()

	requireAcceptanceEnabled(t)
	env := requireEnv(t, envBaseURL, envUsername, envPassword, envOrganizationID)

	organizationID, err := strconv.ParseInt(env[envOrganizationID], 10, 64)
	if err != nil {
		t.Fatalf("invalid %s: %v", envOrganizationID, err)
	}

	awxClient := mustClient(t, env)
	ctx := context.Background()
	resourceName := fmt.Sprintf("tf-awx-acceptance-team-lookup-%d", time.Now().UnixNano())

	t.Logf("creating lookup fixture team %q in organization %d", resourceName, organizationID)
	created, err := awxClient.CreateObject(ctx, "/api/v2/teams/", map[string]any{
		"name":         resourceName,
		"description":  "terraform provider data-source lookup acceptance test",
		"organization": organizationID,
	})
	if err != nil {
		t.Fatalf("failed to create team: %v", err)
	}

	teamID, err := parseID(created["id"])
	if err != nil {
		t.Fatalf("failed to parse created team id: %v", err)
	}
	defer func() {
		t.Logf("cleanup: deleting lookup fixture team id=%d", teamID)
		_ = awxClient.DeleteObject(ctx, "/api/v2/teams/{id}/", strconv.FormatInt(teamID, 10))
	}()

	t.Logf("looking up team by id=%d", teamID)
	byID, err := awxClient.GetObject(ctx, "/api/v2/teams/{id}/", strconv.FormatInt(teamID, 10))
	if err != nil {
		t.Fatalf("lookup by id failed: %v", err)
	}
	byIDValue, err := parseID(byID["id"])
	if err != nil {
		t.Fatalf("failed to parse id from id-lookup: %v", err)
	}
	if byIDValue != teamID {
		t.Fatalf("unexpected id from id-lookup: got=%d want=%d", byIDValue, teamID)
	}

	t.Logf("looking up team by name=%q", resourceName)
	byName, err := awxClient.FindByField(ctx, "/api/v2/teams/", "name", resourceName)
	if err != nil {
		t.Fatalf("lookup by name failed: %v", err)
	}
	if len(byName) != 1 {
		t.Fatalf("expected exactly one result for unique name lookup, got=%d", len(byName))
	}
	byNameValue, err := parseID(byName[0]["id"])
	if err != nil {
		t.Fatalf("failed to parse id from name-lookup: %v", err)
	}
	if byNameValue != teamID {
		t.Fatalf("unexpected id from name-lookup: got=%d want=%d", byNameValue, teamID)
	}
}

func requireAcceptanceEnabled(t *testing.T) {
	t.Helper()
	if os.Getenv(envAcceptance) != "1" {
		t.Skipf("Acceptance tests are opt-in. Set %s=1 plus AWX connection env vars to run.", envAcceptance)
	}
}

func requireEnv(t *testing.T, keys ...string) map[string]string {
	t.Helper()
	values := make(map[string]string, len(keys))
	missing := make([]string, 0)
	for _, key := range keys {
		value := os.Getenv(key)
		if value == "" {
			missing = append(missing, key)
			continue
		}
		values[key] = value
	}

	if len(missing) > 0 {
		t.Skipf("Missing required acceptance env vars: %v", missing)
	}
	return values
}

func mustClient(t *testing.T, env map[string]string) *client.Client {
	t.Helper()
	awxClient, err := client.New(client.Config{
		BaseURL:             env[envBaseURL],
		Username:            env[envUsername],
		Password:            env[envPassword],
		RetryMaxAttempts:    3,
		RetryInitialBackoff: 200 * time.Millisecond,
		Timeout:             30 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return awxClient
}

func parseID(raw any) (int64, error) {
	switch value := raw.(type) {
	case float64:
		return int64(value), nil
	case int64:
		return value, nil
	case int:
		return int64(value), nil
	case string:
		return strconv.ParseInt(value, 10, 64)
	default:
		return 0, fmt.Errorf("unsupported id type %T", raw)
	}
}
