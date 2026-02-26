package main

import (
	"context"
	"log"

	"github.com/damien/terraform-provider-awx-iwd/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const providerAddress = "registry.terraform.io/iwd/awx"

func main() {
	err := providerserver.Serve(context.Background(), provider.New("dev"), providerserver.ServeOpts{
		Address: providerAddress,
	})
	if err != nil {
		log.Fatalf("failed to serve provider: %v", err)
	}
}
