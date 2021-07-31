// Copyright 2021 The Rode Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"github.com/alexashley/terraform-provider-rode/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"log"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// populated by goreleaser
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
