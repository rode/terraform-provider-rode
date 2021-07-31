package main

import (
	"context"
	_ "embed"
	"flag"
	"github.com/alexashley/terraform-provider-rode/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"log"
)

//go:embed version
var version string

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set flag for debugger support")
	flag.Parse()

	options := &plugin.ServeOpts{
		ProviderFunc: provider.New(version),
	}

	if debug {
		if err := plugin.Debug(context.Background(), "registry.terraform.io/rode/rode", options); err != nil {
			log.Fatal(err.Error())
		}

		return
	}

	plugin.Serve(options)
}
