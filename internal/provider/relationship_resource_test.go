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
	if _, ok := resp.Schema.Attributes["child_id"]; ok {
		t.Fatalf("did not expect survey spec schema to include child_id")
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

	if _, ok := resp.Schema.Attributes["child_id"]; !ok {
		t.Fatalf("expected association schema to include child_id")
	}
	if _, ok := resp.Schema.Attributes["spec"]; ok {
		t.Fatalf("did not expect association schema to include spec")
	}
}
