package provider

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/damien/terraform-awx-provider/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestResolveObjectDataSourceTargetByID(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{
		Name:           "projects",
		CollectionPath: "/api/v2/projects/",
		DetailPath:     "/api/v2/projects/{id}/",
	}

	var called bool
	target, diags := resolveObjectDataSourceTarget(context.Background(), &fakeObjectLookupClient{
		getObjectFn: func(_ context.Context, detailPath string, id int64) (map[string]any, error) {
			called = true
			if detailPath != object.DetailPath {
				t.Fatalf("unexpected detail path: got=%q want=%q", detailPath, object.DetailPath)
			}
			if id != 42 {
				t.Fatalf("unexpected id: got=%d want=42", id)
			}
			return map[string]any{"id": 42, "name": "demo"}, nil
		},
	}, object, dataSourceLookupInput{
		ID:           types.StringValue("42"),
		HasNameField: true,
	})
	if diags.HasError() {
		t.Fatalf("expected no diagnostics errors, got: %v", diags)
	}
	if !called {
		t.Fatalf("expected GetObject to be called")
	}
	if got := fmt.Sprintf("%v", target["name"]); got != "demo" {
		t.Fatalf("unexpected target name: got=%q want=%q", got, "demo")
	}
}

func TestResolveObjectDataSourceTargetByIDPrefersIDWhenNameAlsoSet(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{
		Name:           "projects",
		CollectionPath: "/api/v2/projects/",
		DetailPath:     "/api/v2/projects/{id}/",
	}

	var getCalled bool
	target, diags := resolveObjectDataSourceTarget(context.Background(), &fakeObjectLookupClient{
		getObjectFn: func(_ context.Context, detailPath string, id int64) (map[string]any, error) {
			getCalled = true
			if detailPath != object.DetailPath {
				t.Fatalf("unexpected detail path: got=%q want=%q", detailPath, object.DetailPath)
			}
			if id != 9 {
				t.Fatalf("unexpected id: got=%d want=9", id)
			}
			return map[string]any{"id": 9, "name": "from-id"}, nil
		},
		findByFieldFn: func(_ context.Context, _, _, _ string) ([]map[string]any, error) {
			t.Fatalf("expected FindByField not to be called when ID is provided")
			return nil, nil
		},
	}, object, dataSourceLookupInput{
		ID:           types.StringValue("9"),
		Name:         types.StringValue("from-name"),
		HasNameField: true,
	})

	if diags.HasError() {
		t.Fatalf("expected no diagnostics errors, got: %v", diags)
	}
	if !getCalled {
		t.Fatalf("expected GetObject to be called")
	}
	if got := fmt.Sprintf("%v", target["name"]); got != "from-id" {
		t.Fatalf("unexpected lookup target: got=%q want=%q", got, "from-id")
	}
}

func TestResolveObjectDataSourceTargetByNameSingleMatch(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{
		Name:           "projects",
		CollectionPath: "/api/v2/projects/",
	}

	var called bool
	target, diags := resolveObjectDataSourceTarget(context.Background(), &fakeObjectLookupClient{
		findByFieldFn: func(_ context.Context, collectionPath, field, target string) ([]map[string]any, error) {
			called = true
			if collectionPath != object.CollectionPath {
				t.Fatalf("unexpected collection path: got=%q want=%q", collectionPath, object.CollectionPath)
			}
			if field != "name" {
				t.Fatalf("unexpected field: got=%q want=%q", field, "name")
			}
			if target != "demo" {
				t.Fatalf("unexpected target value: got=%q want=%q", target, "demo")
			}
			return []map[string]any{{"id": 7, "name": "demo"}}, nil
		},
	}, object, dataSourceLookupInput{
		Name:         types.StringValue("demo"),
		HasNameField: true,
	})
	if diags.HasError() {
		t.Fatalf("expected no diagnostics errors, got: %v", diags)
	}
	if !called {
		t.Fatalf("expected FindByField to be called")
	}
	if got := fmt.Sprintf("%v", target["id"]); got != "7" {
		t.Fatalf("unexpected target id: got=%q want=%q", got, "7")
	}
}

func TestResolveObjectDataSourceTargetNoMatch(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{Name: "projects", CollectionPath: "/api/v2/projects/"}
	_, diags := resolveObjectDataSourceTarget(context.Background(), &fakeObjectLookupClient{
		findByFieldFn: func(_ context.Context, _, _, _ string) ([]map[string]any, error) {
			return []map[string]any{}, nil
		},
	}, object, dataSourceLookupInput{
		Name:         types.StringValue("missing"),
		HasNameField: true,
	})

	if !diags.HasError() {
		t.Fatalf("expected error diagnostics")
	}
	if !hasDiagSummary(diags, "AWX object not found") {
		t.Fatalf("expected AWX object not found diagnostic, got: %v", diags)
	}
}

func TestResolveObjectDataSourceTargetAmbiguous(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{Name: "projects", CollectionPath: "/api/v2/projects/"}
	_, diags := resolveObjectDataSourceTarget(context.Background(), &fakeObjectLookupClient{
		findByFieldFn: func(_ context.Context, _, _, _ string) ([]map[string]any, error) {
			return []map[string]any{
				{"id": 1, "name": "dupe"},
				{"id": 2, "name": "dupe"},
			}, nil
		},
	}, object, dataSourceLookupInput{
		Name:         types.StringValue("dupe"),
		HasNameField: true,
	})

	if !diags.HasError() {
		t.Fatalf("expected error diagnostics")
	}
	if !hasDiagSummary(diags, "Ambiguous AWX object lookup") {
		t.Fatalf("expected ambiguous lookup diagnostic, got: %v", diags)
	}
}

func TestResolveObjectDataSourceTargetMissingLookupInput(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{Name: "projects"}
	_, diags := resolveObjectDataSourceTarget(context.Background(), &fakeObjectLookupClient{}, object, dataSourceLookupInput{
		HasNameField: true,
	})

	if !diags.HasError() {
		t.Fatalf("expected error diagnostics")
	}
	if !hasDiagSummary(diags, "Missing lookup input") {
		t.Fatalf("expected missing lookup input diagnostic, got: %v", diags)
	}
}

func TestResolveObjectDataSourceTargetInvalidID(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{Name: "projects", DetailPath: "/api/v2/projects/{id}/"}
	_, diags := resolveObjectDataSourceTarget(context.Background(), &fakeObjectLookupClient{}, object, dataSourceLookupInput{
		ID: types.StringValue("not-a-number"),
	})

	if !diags.HasError() {
		t.Fatalf("expected error diagnostics")
	}
	if !hasDiagSummary(diags, "Invalid AWX object ID") {
		t.Fatalf("expected invalid id diagnostic, got: %v", diags)
	}
}

func TestResolveObjectDataSourceTargetLookupError(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{Name: "projects", CollectionPath: "/api/v2/projects/"}
	_, diags := resolveObjectDataSourceTarget(context.Background(), &fakeObjectLookupClient{
		findByFieldFn: func(_ context.Context, _, _, _ string) ([]map[string]any, error) {
			return nil, errors.New("boom")
		},
	}, object, dataSourceLookupInput{
		Name:         types.StringValue("demo"),
		HasNameField: true,
	})

	if !diags.HasError() {
		t.Fatalf("expected error diagnostics")
	}
	if !hasDiagSummary(diags, "Failed to query AWX object by name") {
		t.Fatalf("expected lookup failure diagnostic, got: %v", diags)
	}
}

func TestResolveObjectDataSourceTargetGetObjectError(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{Name: "projects", DetailPath: "/api/v2/projects/{id}/"}
	_, diags := resolveObjectDataSourceTarget(context.Background(), &fakeObjectLookupClient{
		getObjectFn: func(_ context.Context, _ string, _ int64) (map[string]any, error) {
			return nil, errors.New("boom")
		},
	}, object, dataSourceLookupInput{
		ID: types.StringValue("12"),
	})

	if !diags.HasError() {
		t.Fatalf("expected error diagnostics")
	}
	if !hasDiagSummary(diags, "Failed to query AWX object") {
		t.Fatalf("expected object query failure diagnostic, got: %v", diags)
	}
}

func TestResolveObjectDataSourceTargetNameIgnoredWhenFieldUnavailable(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{Name: "projects"}
	_, diags := resolveObjectDataSourceTarget(context.Background(), &fakeObjectLookupClient{}, object, dataSourceLookupInput{
		Name:         types.StringValue("demo"),
		HasNameField: false,
	})

	if !diags.HasError() {
		t.Fatalf("expected error diagnostics")
	}
	if !hasDiagSummary(diags, "Missing lookup input") {
		t.Fatalf("expected missing lookup input diagnostic, got: %v", diags)
	}
}

type fakeObjectLookupClient struct {
	getObjectFn   func(context.Context, string, int64) (map[string]any, error)
	findByFieldFn func(context.Context, string, string, string) ([]map[string]any, error)
}

func (f *fakeObjectLookupClient) GetObject(ctx context.Context, detailPath string, id int64) (map[string]any, error) {
	if f.getObjectFn == nil {
		return nil, errors.New("unexpected GetObject call")
	}
	return f.getObjectFn(ctx, detailPath, id)
}

func (f *fakeObjectLookupClient) FindByField(ctx context.Context, endpointPath string, field string, target string) ([]map[string]any, error) {
	if f.findByFieldFn == nil {
		return nil, errors.New("unexpected FindByField call")
	}
	return f.findByFieldFn(ctx, endpointPath, field, target)
}

func hasDiagSummary(diags diag.Diagnostics, summary string) bool {
	for _, diagnostic := range diags {
		if diagnostic.Summary() == summary {
			return true
		}
	}
	return false
}
