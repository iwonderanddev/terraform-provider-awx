package provider

import (
	"context"
	"testing"

	"github.com/damien/terraform-awx-provider/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestRelationshipResourceSchemaSurveySpec(t *testing.T) {
	t.Parallel()

	resourceInstance := NewRelationshipResource(manifest.Relationship{
		Name:         "job_template_survey_spec",
		ResourceName: "awx_job_template_survey_spec",
		ParentObject: "job_templates",
		ChildObject:  "survey_spec",
		Path:         "/api/v2/job_templates/{id}/survey_spec/",
	})

	resp := resource.SchemaResponse{}
	resourceInstance.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	if _, ok := resp.Schema.Attributes["spec"]; !ok {
		t.Fatalf("expected survey spec schema to include spec attribute")
	}
	if _, ok := resp.Schema.Attributes["job_template_id"]; !ok {
		t.Fatalf("expected survey spec schema to include job_template_id attribute")
	}
	if _, ok := resp.Schema.Attributes["parent_id"]; ok {
		t.Fatalf("did not expect survey spec schema to include parent_id attribute")
	}
}

func TestRelationshipResourceSchemaAssociation(t *testing.T) {
	t.Parallel()

	resourceInstance := NewRelationshipResource(manifest.Relationship{
		Name:         "team_user_association",
		ResourceName: "awx_team_user_association",
		ParentObject: "teams",
		ChildObject:  "users",
		Path:         "/api/v2/teams/{id}/users/",
	})

	resp := resource.SchemaResponse{}
	resourceInstance.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	if _, ok := resp.Schema.Attributes["team_id"]; !ok {
		t.Fatalf("expected association schema to include team_id")
	}
	if _, ok := resp.Schema.Attributes["user_id"]; !ok {
		t.Fatalf("expected association schema to include user_id")
	}
	if _, ok := resp.Schema.Attributes["parent_id"]; ok {
		t.Fatalf("did not expect association schema to include parent_id")
	}
	if _, ok := resp.Schema.Attributes["child_id"]; ok {
		t.Fatalf("did not expect association schema to include child_id")
	}
	if _, ok := resp.Schema.Attributes["spec"]; ok {
		t.Fatalf("did not expect association schema to include spec")
	}
}

func TestParseCompositeRelationshipImportID(t *testing.T) {
	t.Parallel()

	parentID, childID, err := parseCompositeRelationshipImportID("12:34")
	if err != nil {
		t.Fatalf("expected valid composite import ID, got error: %v", err)
	}
	if parentID != 12 || childID != 34 {
		t.Fatalf("unexpected parsed IDs: got=%d:%d want=12:34", parentID, childID)
	}

	if _, _, err := parseCompositeRelationshipImportID("12"); err == nil {
		t.Fatalf("expected malformed composite import ID to fail")
	}
}

func TestParseSurveySpecImportID(t *testing.T) {
	t.Parallel()

	parentID, err := parseSurveySpecImportID("12")
	if err != nil {
		t.Fatalf("expected valid survey-spec import ID, got error: %v", err)
	}
	if parentID != 12 {
		t.Fatalf("unexpected parent ID: got=%d want=%d", parentID, 12)
	}

	if _, err := parseSurveySpecImportID("12:34"); err == nil {
		t.Fatalf("expected malformed survey-spec import ID to fail")
	}
}
