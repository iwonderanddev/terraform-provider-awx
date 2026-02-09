package main

import (
	"testing"

	"github.com/damien/terraform-awx-provider/internal/manifest"
)

func TestBuildReportExcludesRuntimeDataSourcesFromManagedCoverage(t *testing.T) {
	t.Parallel()

	objects := []manifest.ManagedObject{
		{
			Name:             "projects",
			ResourceName:     "awx_project",
			DataSourceName:   "awx_project",
			ResourceEligible: true,
			DataSourceElig:   true,
			RuntimeExcluded:  false,
		},
		{
			Name:             "jobs",
			ResourceName:     "awx_job",
			DataSourceName:   "awx_job",
			ResourceEligible: false,
			DataSourceElig:   true,
			RuntimeExcluded:  true,
		},
	}

	report := buildReport("external/awx-openapi/schema.json", objects, nil, map[string]manifest.RuntimeExclusion{
		"jobs": {Object: "jobs", Reason: "runtime"},
	})

	if report.DataSourceEligible != 1 {
		t.Fatalf("unexpected data source eligible count: got=%d want=1", report.DataSourceEligible)
	}
	if report.RuntimeExcluded != 1 {
		t.Fatalf("unexpected runtime excluded count: got=%d want=1", report.RuntimeExcluded)
	}
	if len(report.ManagedDataSourceObjects) != 1 {
		t.Fatalf("unexpected managed data source object count: got=%d want=1", len(report.ManagedDataSourceObjects))
	}
	if report.ManagedDataSourceObjects[0] != "awx_project" {
		t.Fatalf("unexpected managed data source name: got=%q want=%q", report.ManagedDataSourceObjects[0], "awx_project")
	}
}
