package main

import (
	"context"
	"log"

	"github.com/damien/terraform-awx-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const providerAddress = "registry.terraform.io/damien/awx"

func main() {
	err := providerserver.Serve(context.Background(), provider.New("dev"), providerserver.ServeOpts{
		Address: providerAddress,
	})
	if err != nil {
		log.Fatalf("failed to serve provider: %v", err)
	}
}
