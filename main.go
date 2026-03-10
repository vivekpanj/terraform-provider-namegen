package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
)

var (
	version string = "dev"
)

func main() {

	// Create the provider server
	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(New(version)()),
	}

	// Serve the provider
	tf6server.Serve("registry.terraform.io/local/namegen", providers[0])
}