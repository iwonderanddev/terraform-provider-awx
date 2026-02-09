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

	"github.com/damien/terraform-awx-provider/internal/client"
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
		_ = awxClient.DeleteObject(ctx, "/api/v2/teams/{id}/", teamID)
	}()

	if _, err := awxClient.UpdateObject(ctx, "/api/v2/teams/{id}/", teamID, map[string]any{"description": "updated by acceptance test"}); err != nil {
		t.Fatalf("failed to update team: %v", err)
	}

	fetched, err := awxClient.GetObject(ctx, "/api/v2/teams/{id}/", teamID)
	if err != nil {
		t.Fatalf("failed to read team: %v", err)
	}
	if got := fmt.Sprintf("%v", fetched["name"]); got != resourceName {
		t.Fatalf("unexpected team name: got=%q want=%q", got, resourceName)
	}

	importID := strconv.FormatInt(teamID, 10)
	if _, err := strconv.ParseInt(importID, 10, 64); err != nil {
		t.Fatalf("expected numeric import id, got %q", importID)
	}

	if err := awxClient.DeleteObject(ctx, "/api/v2/teams/{id}/", teamID); err != nil {
		t.Fatalf("failed to delete team: %v", err)
	}

	_, err = awxClient.GetObject(ctx, "/api/v2/teams/{id}/", teamID)
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

	existedBefore, err := awxClient.RelationshipExists(ctx, relationshipPath, teamID, userID)
	if err != nil {
		t.Fatalf("failed pre-check for relationship existence: %v", err)
	}

	if err := awxClient.Associate(ctx, relationshipPath, teamID, userID); err != nil {
		t.Fatalf("failed to associate team and user: %v", err)
	}

	existsAfterCreate, err := awxClient.RelationshipExists(ctx, relationshipPath, teamID, userID)
	if err != nil {
		t.Fatalf("failed relationship read after create: %v", err)
	}
	if !existsAfterCreate {
		t.Fatalf("expected relationship to exist after associate")
	}

	compositeID := fmt.Sprintf("%d:%d", teamID, userID)
	if !compositeIDPattern.MatchString(compositeID) {
		t.Fatalf("invalid composite import id: %q", compositeID)
	}

	if err := awxClient.Disassociate(ctx, relationshipPath, teamID, userID); err != nil {
		t.Fatalf("failed to disassociate team and user: %v", err)
	}

	existsAfterDelete, err := awxClient.RelationshipExists(ctx, relationshipPath, teamID, userID)
	if err != nil {
		t.Fatalf("failed relationship read after delete: %v", err)
	}
	if existsAfterDelete {
		t.Fatalf("expected relationship to be absent after disassociate")
	}

	if existedBefore {
		if err := awxClient.Associate(ctx, relationshipPath, teamID, userID); err != nil {
			t.Fatalf("failed to restore original relationship state: %v", err)
		}
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
