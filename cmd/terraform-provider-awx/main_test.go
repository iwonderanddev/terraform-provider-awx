package main

import (
	"os"
	"strings"
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty defaults to dev", input: "", want: "dev"},
		{name: "whitespace defaults to dev", input: "  ", want: "dev"},
		{name: "trimmed semver", input: " v0.1.0 ", want: "v0.1.0"},
		{name: "prerelease allowed", input: "v0.1.1-rc.1", want: "v0.1.1-rc.1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeVersion(tc.input)
			if got != tc.want {
				t.Fatalf("normalizeVersion(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestMakeBuildTargetsTerraformProviderAwxBinaryName(t *testing.T) {
	raw, err := os.ReadFile("../../Makefile")
	if err != nil {
		t.Fatalf("failed to read Makefile: %v", err)
	}

	content := string(raw)
	if !strings.Contains(content, "-o dist/terraform-provider-awx ./cmd/terraform-provider-awx") {
		t.Fatalf("expected make build to emit dist/terraform-provider-awx, got:\n%s", content)
	}
	if strings.Contains(content, "-o dist/terraform-provider-awx-iwd ./cmd/terraform-provider-awx") {
		t.Fatalf("unexpected awx-iwd output name in make build target")
	}
}
