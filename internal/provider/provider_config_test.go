package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	providerframework "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestProviderSchemaUsesHostnameAttribute(t *testing.T) {
	t.Parallel()

	providerInstance := New("test")().(*awxProvider) //nolint:forcetypeassert
	resp := providerframework.SchemaResponse{}
	providerInstance.Schema(context.Background(), providerframework.SchemaRequest{}, &resp)

	if _, ok := resp.Schema.Attributes["hostname"]; !ok {
		t.Fatalf("expected provider schema to expose hostname attribute")
	}
	if _, ok := resp.Schema.Attributes["base_url"]; ok {
		t.Fatalf("expected provider schema to no longer expose base_url attribute")
	}
}

func TestValidateConfigRejectsInvalidHostname(t *testing.T) {
	t.Parallel()

	config := providerModel{
		Hostname: types.StringValue("awx.example.invalid"),
	}

	var diags diag.Diagnostics
	validateConfig(config, &diags)

	if !hasDiagnosticSummary(diags, "Invalid AWX hostname") {
		t.Fatalf("expected Invalid AWX hostname diagnostic, got: %#v", diags)
	}
	if !hasDiagnosticDetailContaining(diags, "hostname must use http or https.") {
		t.Fatalf("expected hostname scheme diagnostic detail, got: %#v", diags)
	}
}

func hasDiagnosticSummary(diags diag.Diagnostics, summary string) bool {
	for _, diagnostic := range diags {
		if diagnostic.Summary() == summary {
			return true
		}
	}
	return false
}

func hasDiagnosticDetailContaining(diags diag.Diagnostics, substring string) bool {
	for _, diagnostic := range diags {
		if strings.Contains(diagnostic.Detail(), substring) {
			return true
		}
	}
	return false
}
