package main

import (
	"context"
	"flag"
	"log"

	"github.com/chatbotkit/terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// version is set via -ldflags at build time
var version string = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/chatbotkit/chatbotkit",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
