package main

import (
	"context"
	"log"
	"strings"

	"github.com/damien/terraform-provider-awx-iwd/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const providerAddress = "registry.terraform.io/iwd/awx"

// version is injected at build time with -ldflags "-X main.version=<semver>".
var version = "dev"

func main() {
	err := providerserver.Serve(context.Background(), provider.New(normalizeVersion(version)), providerserver.ServeOpts{
		Address: providerAddress,
	})
	if err != nil {
		log.Fatalf("failed to serve provider: %v", err)
	}
}

func normalizeVersion(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "dev"
	}
	return trimmed
}
